package loms

import (
	"context"
	"route256/loms/internal/pkg/logger"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) OrderPayed(ctx context.Context, req *loms_v1.OrderPayedRequest) (*emptypb.Empty, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.impl.OrderPayed(ctx, req.Orderid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
