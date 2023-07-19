package loms

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
	"route256/loms/internal/pkg/logger"
	"route256/loms/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) Stocks(ctx context.Context, sku uint32) ([]model.StockItem, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/loms/stocks")
	defer span.Finish()

	span.SetTag("sku", sku)

	logger.Info("listing warehouses for sku: %d", sku)

	stockItem, err := s.ordersRepository.Stocks(ctx, sku)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, fmt.Errorf("stock response: %w", err))
	}

	return stockItem, nil
}
