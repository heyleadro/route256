package checkout

import (
	"context"
	"fmt"
	"route256/checkout/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

/*
DeleteFromCart updates UserCart, may delete Cart entry if Count is 0
called with checkout(localhost):8080/deleteFromCart
*/
func (s *Service) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/checkout/deletefromcart")
	defer span.Finish()

	span.SetTag("user_id", user)
	span.LogKV("sku", sku)
	span.LogKV("count", count)

	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTx context.Context) error {

		cartCount, err := s.cartRepo.GetUserSkuItemCount(ctx, user, sku)
		if err != nil {
			return fmt.Errorf("get item count query: %w", err)
		}

		if cartCount < count {
			return ErrCartInsufficient
		}

		err = s.cartRepo.DeleteFromCart(ctx, user, sku, count)
		if err != nil {
			return fmt.Errorf("delete from cart query: %w", err)
		}

		return nil
	})
	if err != nil {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("delete from cart tx: %w", err))
	}

	return nil
}
