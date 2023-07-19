package schema

type StockItem struct {
	WarehouseID int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}

type OrderDescr struct {
	UserID int64  `db:"user_id"`
	Status string `db:"status"`
}

type Item struct {
	SKU   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}

type OrderItem struct {
	WarehouseID int64  `db:"warehouse_id"`
	SKU         uint32 `db:"sku"`
	Count       uint16 `db:"count"`
}
