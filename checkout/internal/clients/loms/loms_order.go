package loms

import (
	"context"
	"fmt"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/model"
)

// const (
// 	orderPath = "createOrder"
// )

// CreateOrder invokes loms:8081/createOrder
func (c *Client) Purchase(ctx context.Context, user int64, items []model.Cart) (int64, error) {
	requestStocks := converter.CreateOrderReqestToLomsRequest(user, items)

	responseStocks, err := c.LomsClient.CreateOrder(ctx, &requestStocks)
	if err != nil {
		return -1, fmt.Errorf("order bad request: %w", err)
	}

	return responseStocks.Orderid, nil
}
