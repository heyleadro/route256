package schema

type Cart struct {
	SKU   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}
