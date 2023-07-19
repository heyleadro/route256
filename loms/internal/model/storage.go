package model

import (
	"errors"
)

const ServiceName = "loms"

type OrderStorage struct { // key is orderID
	Orders map[int64]Order
}

type Order struct {
	Status string
	User   int64
	Items  []UserItems
}

type UserItems struct {
	WarehouseID int64
	SKU         uint32
	Count       uint16
}

type Warehouse struct {
	Warehouses []StockItem
}

type StockItem struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

func (w *Warehouse) UpdateWarehouse(stockItem StockItem) error {
	for idx, v := range w.Warehouses {
		if v.WarehouseID == stockItem.WarehouseID {
			if w.Warehouses[idx].Count < stockItem.Count {
				return ErrNotSuffientStock
			}
			w.Warehouses[idx].Count -= stockItem.Count
			break
		}
	}
	return nil
}

func (w *Warehouse) ReturnReserved(stockItem StockItem) {
	for idx, v := range w.Warehouses {
		if v.WarehouseID == stockItem.WarehouseID {
			w.Warehouses[idx].Count += stockItem.Count
			return
		}
	}
}

func (w *Warehouse) GetWarehouses() ([]StockItem, error) {
	if w.Warehouses == nil {
		return nil, ErrNoWarehouses
	}

	return w.Warehouses, nil
}

// PutOrder updates storage of all orders
func (o *OrderStorage) PutOrder(orderID int64, user int64, userItems []UserItems, status string) {
	o.Orders[orderID] = Order{
		Status: status,
		User:   user,
		Items:  userItems,
	}
}

func (o *OrderStorage) UpdateOrder(orderID int64, status string) {
	orderCopy := o.Orders[orderID]

	orderCopy.Status = status

	o.Orders[orderID] = orderCopy
}

func (o *OrderStorage) GetOrder(orderID int64) (Order, error) {
	order, err := o.Orders[orderID]
	if !err {
		return Order{}, ErrNoOrder
	}

	return order, nil
}

// initializing warehouse for stocks
func NewWarehouse() *Warehouse {
	return &Warehouse{
		Warehouses: []StockItem{
			{WarehouseID: 1, Count: 200},
		},
	}
}

// initializing order storage for create_order, list_order etc.
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		Orders: make(map[int64]Order),
	}
}

type ProducerMessage struct {
	UserID  int64
	OrderID int64
	Status  string
}

var (
	ErrNotSuffientStock = errors.New("no items in warehouses")
	ErrNoOrder          = errors.New("no order for orderID")
	ErrNoWarehouses     = errors.New("no warehouses")
)
