package loms

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/pkg/logger"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
Handle responses to addtocart checkout with vector of stockItems(warehouseID and Count)
for now NOT depending on SKU
*/
func (s *Service) Stocks(ctx context.Context, req *loms_v1.StocksRequest) (*loms_v1.StocksResponse, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	stockItems, err := s.impl.Stocks(ctx, req.Sku)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return converter.StocksToResponse(stockItems), nil
}
