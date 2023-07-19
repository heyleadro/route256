package converter

import (
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"
)

func CreateOrderItemsFromRequest(req *loms_v1.CreateOrderRequest) (user int64, userItems []model.UserItems) {
	user = req.User
	userItems = make([]model.UserItems, 0)
	for _, val := range req.Items {
		userItems = append(userItems, model.UserItems{
			SKU:   val.Sku,
			Count: uint16(val.Count),
		})
	}

	return user, userItems
}
