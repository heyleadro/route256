package loms

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
	"route256/loms/internal/pkg/logger"
	"route256/loms/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

// ListOrder return Order struct
func (s *Service) ListOrder(ctx context.Context, orderID int64) (model.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/loms/cancelorder")
	defer span.Finish()

	span.SetTag("order_id", orderID)

	logger.Info("listing order for orderID: %d", orderID)

	order, err := s.ordersRepository.ListOrder(ctx, orderID)
	if err != nil {
		return model.Order{}, tracer.MarkSpanWithError(ctx, fmt.Errorf("get order response: %w", err))
	}

	return order, nil
}
