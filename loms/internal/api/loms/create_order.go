package loms

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/pkg/logger"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handle returns orderID generated in CreateOrder which updates orderstorage and warehouses
func (s *Service) CreateOrder(ctx context.Context, req *loms_v1.CreateOrderRequest) (*loms_v1.CreateOrderResponse, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	user, userItems := converter.CreateOrderItemsFromRequest(req)

	orderID, err := s.impl.CreateOrder(ctx, user, userItems)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &loms_v1.CreateOrderResponse{Orderid: orderID}, nil
}
