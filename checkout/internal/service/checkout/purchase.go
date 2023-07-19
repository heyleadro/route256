package checkout

import (
	"context"
	"fmt"
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/internal/pkg/tracer"

	"github.com/opentracing/opentracing-go"
)

/*
Purchase invokes createOrder in loms(clients/loms/loms_order) with users items in cart
called with checkout(localhost):8080/purchase
*/
func (s *Service) Purchase(ctx context.Context, user int64) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/checkout/purchase")
	defer span.Finish()

	span.SetTag("user_id", user)

	cartItems, err := s.cartRepo.GetUserItems(ctx, user)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, ErrNoUserCart)
	}

	orderID, err := s.orderCreator.Purchase(ctx, user, cartItems)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, fmt.Errorf("purchase: %w", err))
	}

	logger.Info("User's %d orderID: %d", user, orderID)

	return orderID, nil
}
