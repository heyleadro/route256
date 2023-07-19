package checkout

import (
	"context"
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Purchase(ctx context.Context, req *checkout_v1.PurchaseRequest) (*checkout_v1.PurchaseResponse, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// ctx, fnCancel := context.WithTimeout(ctx, 20*time.Second)
	// defer fnCancel()
	orderID, err := s.impl.Purchase(ctx, req.User)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &checkout_v1.PurchaseResponse{Orderid: orderID}, nil
}
