package postgres

import (
	"context"
	"fmt"
	"route256/checkout/internal/model"
	"route256/checkout/internal/repository/postgres/tx"
	"route256/checkout/internal/repository/schema"

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
	tableNameCart = "cart"
)

func (r *Repository) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableNameCart).
		Columns("user_id", "sku", "count").
		Values(user, sku, count).
		Suffix(`ON CONFLICT ("user_id", "sku") DO UPDATE SET count=cart.count+excluded.count;`)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for add to cart: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec add to cart: %w", err)
	}

	return nil

}

func (r *Repository) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	db := r.provider.GetDB(ctx)

	query := psql.Update(tableNameCart).
		Set("count", sq.Expr("count - ?", count)).
		Where(sq.And{sq.Eq{"user_id": user}, sq.Eq{"sku": sku}})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for delete from cart: %w", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec delete from cart: %w", err)
	}

	return nil
}

func (r *Repository) GetUserSkuItemCount(ctx context.Context, user int64, sku uint32) (uint16, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select("count").
		From(tableNameCart).
		Where(sq.And{sq.Eq{"user_id": user}, sq.Eq{"sku": sku}})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for get item count: %w", err)
	}

	var result int
	if err = db.QueryRow(ctx, rawSQL, args...).Scan(&result); err != nil {
		return 0, fmt.Errorf("select get item count: %w", err)
	}

	return uint16(result), nil
}

func (r *Repository) GetUserItems(ctx context.Context, user int64) ([]model.Cart, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select([]string{"sku", "count"}...).
		From(tableNameCart).
		Where("user_id = ?", user)

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for get items: %w", err)
	}

	var result []schema.Cart
	if err = pgxscan.Select(ctx, db, &result, rawSQL, args...); err != nil {
		return nil, fmt.Errorf("select get items: %w", err)
	}

	return bindToModelCarts(result), nil
}

func bindToModelCarts(in []schema.Cart) []model.Cart {
	result := make([]model.Cart, 0, len(in))
	for _, val := range in {
		result = append(result, bindtoModelCart(val))
	}

	return result
}

func bindtoModelCart(in schema.Cart) model.Cart {
	return model.Cart{
		SKU:   in.SKU,
		Count: in.Count,
	}
}
