package converter

import (
	"route256/checkout/internal/model"
	"route256/checkout/pkg/loms_v1"
)

func StocksFromResponse(resp *loms_v1.StocksResponse) []model.Stock {
	resultResponse := make([]model.Stock, 0)

	for _, stock := range resp.Stocks {
		resultResponse = append(resultResponse, model.Stock{
			WarehouseID: stock.Warehouseid,
			Count:       stock.Count,
		})
	}

	return resultResponse
}
