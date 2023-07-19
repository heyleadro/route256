package loms

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/pkg/logger"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handle responses with order obtained with orderID in ListOrder
func (s *Service) ListOrder(ctx context.Context, req *loms_v1.ListOrderRequest) (*loms_v1.ListOrderResponse, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	order, err := s.impl.ListOrder(ctx, req.Orderid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("orderID: %d status: %s, user: %d", req.Orderid, order.Status, order.User)

	logger.Info("Order's Items: ")
	for _, item := range order.Items {
		logger.Info("SKU: %d, Count: %d", item.SKU, item.Count)
	}

	return converter.ListOrderToResponse(order), nil
}
