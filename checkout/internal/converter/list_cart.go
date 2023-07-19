package converter

import (
	"route256/checkout/internal/model"
	"route256/checkout/pkg/checkout_v1"
)

func ListCartToResponse(resp *model.ListCartDescr) *checkout_v1.ListCartResponse {
	items := make([]*checkout_v1.Item, 0)

	for _, val := range resp.ItemsDescr {
		items = append(items, &checkout_v1.Item{
			Sku:   val.CartDescr.SKU,
			Count: uint32(val.CartDescr.Count),
			Name:  val.UserItemDescr.Name,
			Price: val.UserItemDescr.Price,
		})
	}

	return &checkout_v1.ListCartResponse{
		Items:      items,
		Totalprice: resp.TotalPrice,
	}
}
