package model

const ServiceName = "checkout"

type UserCarts struct { // key is user
	Carts map[int64][]Cart
}

type Cart struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Item struct {
	Name  string
	Price uint32
}

type ItemDescr struct {
	CartDescr     Cart
	UserItemDescr Item
}

type ListCartDescr struct {
	TotalPrice uint32
	ItemsDescr []ItemDescr
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

// GetCarts returns copy of UserCarts
func (c *UserCarts) GetCarts() map[int64][]Cart {
	return c.Carts
}

// UpdateCarts update this UserCart with input
func (c *UserCarts) UpdateCarts(m map[int64][]Cart) {
	c.Carts = m
}

func NewCarts() *UserCarts {
	return &UserCarts{
		Carts: make(map[int64][]Cart),
	}
}
