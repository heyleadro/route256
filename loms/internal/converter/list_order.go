package converter

import (
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"
)

func ListOrderToResponse(order model.Order) *loms_v1.ListOrderResponse {
	lomsItems := make([]*loms_v1.UserItems, 0)

	for _, val := range order.Items {
		lomsItems = append(lomsItems, &loms_v1.UserItems{
			Sku:   val.SKU,
			Count: uint32(val.Count),
		})
	}

	return &loms_v1.ListOrderResponse{
		Status: order.Status,
		User:   order.User,
		Items:  lomsItems,
	}
}
