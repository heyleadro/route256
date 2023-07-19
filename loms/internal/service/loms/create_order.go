package loms

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
	"route256/loms/internal/pkg/logger"
	"route256/loms/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

// CreateOrder invoked on createOrder Handle with request from checkout on purchase
func (s *Service) CreateOrder(ctx context.Context, user int64, userItems []model.UserItems) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/loms/createorder")
	defer span.Finish()

	span.SetTag("user_id", user)
	span.LogKV("user_items", userItems)

	logger.Info("creating order for user: %d", user)

	var orderID int64

	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		toBeOrdered := make([]model.UserItems, 0, len(userItems))
		for _, item := range userItems {
			stocks, err := s.ordersRepository.Stocks(ctxTx, item.SKU)
			if err != nil {
				return fmt.Errorf("stock query db: %w", err)
			}

			var reservedCount uint64
			var reservedStock []model.StockItem

			for _, stock := range stocks {
				currentAmount := uint64(item.Count) - reservedCount // how much more items do we need to reserve from another stock

				if currentAmount == 0 {
					break
				}

				if stock.Count >= currentAmount {
					reservedStock = append(reservedStock, model.StockItem{
						WarehouseID: stock.WarehouseID,
						Count:       currentAmount,
					})
					reservedCount += currentAmount
				} else {
					reservedStock = append(reservedStock, model.StockItem{
						WarehouseID: stock.WarehouseID,
						Count:       stock.Count,
					})
					reservedCount += stock.Count
				}
			}

			if reservedCount != uint64(item.Count) {
				return fmt.Errorf("not enough stock for item: %d", item.SKU)
			}

			for _, reservation := range reservedStock {
				err := s.ordersRepository.ReserveStock(ctxTx, item.SKU, reservation)
				if err != nil {
					return fmt.Errorf("query reserve stock: %w", err)
				}

				toBeOrdered = append(toBeOrdered, model.UserItems{
					WarehouseID: reservation.WarehouseID,
					SKU:         item.SKU,
					Count:       uint16(reservation.Count),
				})
			}
		}

		dbOrderID, err := s.ordersRepository.CreateOrder(ctxTx, user, toBeOrdered)
		if err != nil {
			return fmt.Errorf("query create order: %w", err)
		}

		err = s.ordersRepository.UpdateOrderStatus(ctxTx, dbOrderID, "awaiting payment")
		if err != nil {
			return fmt.Errorf("update order status query: %w", err)
		}

		orderID = dbOrderID
		return nil
	})
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, fmt.Errorf("create order tx: %w", err))
	}

	userID, err := s.ordersRepository.GetOrderUser(ctx, orderID)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, fmt.Errorf("no user: %w", err))
	}
	err = s.Notify(ctx, userID, orderID, "awaiting payment")
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, fmt.Errorf("create order notify: %w", err))
	}

	return orderID, nil
}
