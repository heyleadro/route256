package checkout

import (
	"route256/checkout/internal/service/checkout"
	"route256/checkout/pkg/checkout_v1"
)

type Service struct {
	checkout_v1.UnimplementedCheckoutServer
	impl *checkout.Service
}

func NewCheckoutServer(impl *checkout.Service) checkout_v1.CheckoutServer {
	return &Service{impl: impl}
}
