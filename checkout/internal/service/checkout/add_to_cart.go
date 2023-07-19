package checkout

import (
	"context"
	"fmt"
	"route256/checkout/internal/model"
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

/*
AddToCart invokes stocks in loms(clients/loms/loms_stocks), updates UserCarts
called with checkout(localhost):8080/addToCart
*/
func (s *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/checkout/addtocart")
	defer span.Finish()

	span.SetTag("user_id", user)
	span.LogKV("sku", sku)
	span.LogKV("count", count)
	if count <= 0 {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("count must be ge 0"))
	}

	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {
		stocks, err := s.stockChecker.Stocks(ctx, sku) // request of stock in loms
		if err != nil {
			return tracer.MarkSpanWithError(ctx, fmt.Errorf("get stocks: %w", err))
		}

		logger.Info("stocks: %v", stocks)

		if !isSufficient(uint64(count), stocks) {
			return tracer.MarkSpanWithError(ctx, ErrStockInsufficient)
		}

		err = s.cartRepo.AddToCart(ctx, user, sku, count)
		if err != nil {
			return tracer.MarkSpanWithError(ctx, fmt.Errorf("add to cart query: %w", err))
		}

		return nil
	})
	if err != nil {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("add to cart tx: %w", err))
	}

	return nil
}

func isSufficient(count uint64, stocks []model.Stock) bool {
	stockSupply := uint64(0)
	for _, stock := range stocks {
		stockSupply += stock.Count
		if stockSupply >= count {
			break
		}
	}

	return stockSupply >= count
}
