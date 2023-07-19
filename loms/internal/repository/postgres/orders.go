package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/postgres/tx"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func NewRepo(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

const (
	tableNameOrders        = "orders"
	tableNameUserItems     = "user_items"
	tableNameStockItems    = "stock_items"
	tableNameReservedItems = "reserved_items"
)

func (r *Repository) CreateOrder(ctx context.Context, user int64, userItems []model.UserItems) (int64, error) {
	db := r.provider.GetDB(ctx)

	queryInsertOrder := psql.Insert(tableNameOrders).
		Columns("user_id", "status").
		Values(user, "new").
		Suffix("RETURNING order_id")

	rawSQL, args, err := queryInsertOrder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for insert order: %w", err)
	}

	var orderID int64
	if err := db.QueryRow(ctx, rawSQL, args...).Scan(&orderID); err != nil {
		return 0, fmt.Errorf("error create order: %w", err)
	}

	queryInsertItems := psql.Insert(tableNameUserItems).Columns("sku", "order_id", "count", "warehouse_id")
	for _, item := range userItems {
		queryInsertItems = queryInsertItems.Values(item.SKU, orderID, item.Count, item.WarehouseID)
	}

	rawSQL, args, err = queryInsertItems.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for insert user items: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return 0, fmt.Errorf("exec create order: %w", err)
	}

	return orderID, nil

}

func (r *Repository) Stocks(ctx context.Context, sku uint32) ([]model.StockItem, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select("warehouse_id", "count").
		From(tableNameStockItems).
		Where("sku = ?", sku)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error build query: %w", err)
	}

	var result []schema.StockItem
	if err = pgxscan.Select(ctx, db, &result, rawSQL, args...); err != nil {
		return nil, fmt.Errorf("error select stocks: %w", err)
	}

	return bindToModelStockItems(result), nil
}

func bindToModelStockItems(in []schema.StockItem) []model.StockItem {
	result := make([]model.StockItem, 0, len(in))
	for _, val := range in {
		result = append(result, bindToModelStockItem(val))
	}
	return result

}

func bindToModelStockItem(in schema.StockItem) model.StockItem {
	return model.StockItem{
		WarehouseID: in.WarehouseID,
		Count:       in.Count,
	}
}

func (r *Repository) ReserveStock(ctx context.Context, sku uint32, stockItem model.StockItem) error {
	db := r.provider.GetDB(ctx)

	queryReserve := psql.Insert(tableNameReservedItems).
		Columns("sku", "warehouse_id", "count").
		Values(sku, stockItem.WarehouseID, stockItem.Count).
		Suffix(`ON CONFLICT ("sku", "warehouse_id") DO UPDATE SET count=reserved_items.count+excluded.count;`)

	rawSQL, args, err := queryReserve.ToSql()
	if err != nil {
		return fmt.Errorf("build query for insert stocks reservation: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec insert stocks in reservation: %w", err)
	}

	querySubstractStocks := psql.Update(tableNameStockItems).
		Set("count", sq.Expr("count - ?", stockItem.Count)).
		Where(sq.And{sq.Eq{"warehouse_id": stockItem.WarehouseID}, sq.Eq{"sku": sku}})

	rawSQL, args, err = querySubstractStocks.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update stocks: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update stocks: %w", err)
	}

	return nil
}

func (r *Repository) GetOrderStatus(ctx context.Context, orderID int64) (string, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select("status").
		From(tableNameOrders).
		Where("order_id = ?", orderID)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return "", fmt.Errorf("build query for status: %w", err)
	}

	var result string
	if err = db.QueryRow(ctx, rawSQL, args...).Scan(&result); err != nil {
		return "", fmt.Errorf("select order status: %w", err)
	}

	return result, nil
}

func (r *Repository) GetOrderUser(ctx context.Context, orderID int64) (int64, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select("user_id").
		From(tableNameOrders).
		Where("order_id = ?", orderID)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for status: %w", err)
	}

	var result int64
	if err = db.QueryRow(ctx, rawSQL, args...).Scan(&result); err != nil {
		return 0, fmt.Errorf("select order status: %w", err)
	}

	return result, nil
}

func (r *Repository) ListOrder(ctx context.Context, orderID int64) (model.Order, error) {
	db := r.provider.GetDB(ctx)

	queryOrders := psql.Select("user_id", "status").
		From(tableNameOrders).
		Where("order_id = ?", orderID)

	rawSQL, args, err := queryOrders.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("build query for orders: %w", err)
	}

	var resultOrder schema.OrderDescr
	if err = pgxscan.Get(ctx, db, &resultOrder, rawSQL, args...); err != nil {
		return model.Order{}, fmt.Errorf("select orders: %w", err)
	}

	queryItems := psql.Select([]string{"sku", "count"}...).
		From(tableNameUserItems).
		Where("order_id = ?", orderID)

	rawSQL, args, err = queryItems.ToSql()
	if err != nil {
		return model.Order{}, fmt.Errorf("build query for items: %w", err)
	}

	var resultItems []schema.Item
	if err = pgxscan.Select(ctx, db, &resultItems, rawSQL, args...); err != nil {
		return model.Order{}, fmt.Errorf("select items: %w", err)
	}

	return bindToModelOrder(resultOrder, resultItems), nil
}

func bindToModelOrder(order schema.OrderDescr, items []schema.Item) model.Order {
	userItems := bindToModelUserItems(items)

	return model.Order{
		Status: order.Status,
		User:   order.UserID,
		Items:  userItems,
	}
}

func bindToModelUserItems(items []schema.Item) []model.UserItems {
	result := make([]model.UserItems, 0, len(items))
	for _, val := range items {
		result = append(result, bindToModelUserItem(val))
	}
	return result
}

func bindToModelUserItem(item schema.Item) model.UserItems {
	return model.UserItems{
		SKU:   item.SKU,
		Count: item.Count,
	}
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	db := r.provider.GetDB(ctx)
	query := psql.Update(tableNameOrders).
		Set("status", status).
		Where(sq.Eq{"order_id": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update status: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update status: %w", err)
	}

	return nil
}

func (r *Repository) ReturnStock(ctx context.Context, userItems model.UserItems) error {
	db := r.provider.GetDB(ctx)

	queryReturn := psql.Update(tableNameReservedItems).
		Set("count", sq.Expr("count - ?", userItems.Count)).
		Where(sq.And{sq.Eq{"warehouse_id": userItems.WarehouseID}, sq.Eq{"sku": userItems.SKU}})

	rawSQL, args, err := queryReturn.ToSql()
	if err != nil {
		return fmt.Errorf("build query for return stocks reservation: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec return stocks in reservation: %w", err)
	}

	queryAddStocks := psql.Update(tableNameStockItems).
		Set("count", sq.Expr("count + ?", userItems.Count)).
		Where(sq.And{sq.Eq{"warehouse_id": userItems.WarehouseID}, sq.Eq{"sku": userItems.SKU}})

	rawSQL, args, err = queryAddStocks.ToSql()
	if err != nil {
		return fmt.Errorf("build query for update stocks: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec update stocks: %w", err)
	}

	return nil
}

func (r *Repository) FreeReservedStock(ctx context.Context, userItems model.UserItems) error {
	db := r.provider.GetDB(ctx)

	queryReturn := psql.Update(tableNameReservedItems).
		Set("count", sq.Expr("count - ?", userItems.Count)).
		Where(sq.And{sq.Eq{"warehouse_id": userItems.WarehouseID}, sq.Eq{"sku": userItems.SKU}})

	rawSQL, args, err := queryReturn.ToSql()
	if err != nil {
		return fmt.Errorf("build query for return stocks reservation: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec return stocks in reservation: %w", err)
	}

	return nil
}

func (r *Repository) GetOrderItems(ctx context.Context, orderID int64) ([]model.UserItems, error) {
	db := r.provider.GetDB(ctx)

	queryItems := psql.Select([]string{"sku", "count", "warehouse_id"}...).
		From(tableNameUserItems).
		Where("order_id = ?", orderID)

	rawSQL, args, err := queryItems.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for items: %w", err)
	}

	var resultItems []schema.OrderItem
	if err = pgxscan.Select(ctx, db, &resultItems, rawSQL, args...); err != nil {
		return nil, fmt.Errorf("select items: %w", err)
	}

	return bindToModelItem(resultItems), nil
}

func bindToModelItem(in []schema.OrderItem) []model.UserItems {
	result := make([]model.UserItems, 0, len(in))

	for _, val := range in {
		result = append(result, model.UserItems{
			WarehouseID: val.WarehouseID,
			SKU:         val.SKU,
			Count:       val.Count,
		})
	}

	return result
}
