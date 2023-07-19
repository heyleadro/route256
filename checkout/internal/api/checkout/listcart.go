package checkout

import (
	"context"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) ListCart(ctx context.Context, req *checkout_v1.ListCartRequest) (*checkout_v1.ListCartResponse, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	response, err := s.impl.ListCart(ctx, req.User)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return converter.ListCartToResponse(response), nil
}
