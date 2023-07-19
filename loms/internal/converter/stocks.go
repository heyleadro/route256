package converter

import (
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"
)

func StocksToResponse(stockItems []model.StockItem) *loms_v1.StocksResponse {
	stocks := make([]*loms_v1.StockItem, 0)

	for _, val := range stockItems {
		stocks = append(stocks, &loms_v1.StockItem{
			Warehouseid: val.WarehouseID,
			Count:       val.Count,
		})
	}
	return &loms_v1.StocksResponse{Stocks: stocks}
}
