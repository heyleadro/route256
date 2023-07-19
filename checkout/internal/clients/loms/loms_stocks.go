package loms

import (
	"context"
	"fmt"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/loms_v1"
)

// const (
// 	stocksPath = "stocks"
// )

// Stocks invokes loms:8081/stocks
func (c *Client) Stocks(ctx context.Context, sku uint32) ([]model.Stock, error) {
	requestStocks := loms_v1.StocksRequest{Sku: sku}

	responseStocks, err := c.LomsClient.Stocks(ctx, &requestStocks)
	if err != nil {
		return nil, fmt.Errorf("stocks bad request: %w", err)
	}

	return converter.StocksFromResponse(responseStocks), nil
}
