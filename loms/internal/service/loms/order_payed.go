package loms

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/logger"
	"route256/loms/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

// OrderPayed just updates orderstorage with status payed, because "item reservation" is dummy
func (s *Service) OrderPayed(ctx context.Context, orderID int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/loms/orderpayed")
	defer span.Finish()

	span.SetTag("order_id", orderID)

	logger.Info("paying order: %d", orderID)

	status, err := s.ordersRepository.GetOrderStatus(ctx, orderID)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("get order status response: %w", err))
	}

	if status != "awaiting payment" {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("payment status incorrect: %s", status))
	}

	err = s.transactionManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		userItems, err := s.ordersRepository.GetOrderItems(ctxTx, orderID)
		if err != nil {
			return fmt.Errorf("get order items query: %w", err)
		}

		err = s.ordersRepository.UpdateOrderStatus(ctxTx, orderID, "payed")
		if err != nil {
			return fmt.Errorf("update order status query: %w", err)
		}

		for _, item := range userItems {
			err := s.ordersRepository.FreeReservedStock(ctxTx, item)
			if err != nil {
				return fmt.Errorf("free reserved items query: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("order payed tx: %w", err))
	}

	userID, err := s.ordersRepository.GetOrderUser(ctx, orderID)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("no user: %w", err))
	}

	err = s.Notify(ctx, userID, orderID, "payed")
	if err != nil {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("order payed notify: %w", err))
	}

	logger.Info("payed for order: %d, status: %s", orderID, "payed")

	return nil
}
