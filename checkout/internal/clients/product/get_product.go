package product

import (
	"context"
	"fmt"
	"route256/checkout/internal/config"
	"route256/checkout/internal/converter"
	"route256/checkout/internal/model"
	"route256/checkout/pkg/product_v1"
)

// GetProducts makes get_product request to productservices
func (c *Client) GetProducts(ctx context.Context, sku uint32) (model.Item, error) {
	err := c.Limiter.Wait(ctx)
	if err != nil {
		return model.Item{}, fmt.Errorf("rate limiter: %w", err)
	}

	requestProducts := product_v1.GetProductRequest{Token: config.AppConfig.Token, Sku: sku}

	responseProducts, err := c.ProductServiceClient.GetProduct(ctx, &requestProducts)
	if err != nil {
		return model.Item{}, fmt.Errorf("product bad request: %w", err)
	}

	return converter.ItemFromGetProductService(responseProducts), nil
}
