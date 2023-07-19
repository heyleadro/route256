package converter

import (
	"route256/checkout/internal/model"
	"route256/checkout/pkg/loms_v1"
)

func CreateOrderReqestToLomsRequest(user int64, items []model.Cart) loms_v1.CreateOrderRequest {
	lomsItems := make([]*loms_v1.UserItems, 0)

	for _, val := range items {
		lomsItems = append(lomsItems, &loms_v1.UserItems{
			Sku:   val.SKU,
			Count: uint32(val.Count),
		})
	}

	return loms_v1.CreateOrderRequest{User: user, Items: lomsItems}
}
