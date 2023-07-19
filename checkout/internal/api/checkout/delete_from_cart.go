package checkout

import (
	"context"
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) DeleteFromCart(ctx context.Context, req *checkout_v1.DeleteFromCartRequest) (*emptypb.Empty, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.impl.DeleteFromCart(ctx, req.User, req.Sku, uint16(req.Count))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
