package converter

import (
	"route256/checkout/internal/model"
	"route256/checkout/pkg/product_v1"
)

func ItemFromGetProductService(resp *product_v1.GetProductResponse) model.Item {
	return model.Item{
		Name:  resp.Name,
		Price: resp.Price,
	}
}
