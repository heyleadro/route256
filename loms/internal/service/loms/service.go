package loms

import (
	"context"
	"route256/loms/internal/model"
)

// type Warehouser interface {
// 	UpdateWarehouse(stockItem model.StockItem) error
// 	GetWarehouses() ([]model.StockItem, error)
// 	ReturnReserved(stockItem model.StockItem)
// }

// type OrderStorager interface {
// 	PutOrder(orderID int64, user int64, userItems []model.UserItems, status string)
// 	GetOrder(orderID int64) (model.Order, error)
// 	UpdateOrder(orderID int64, status string)
// }

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type OrdersRepository interface {
	CreateOrder(ctx context.Context, user int64, userItems []model.UserItems) (int64, error)
	Stocks(ctx context.Context, sku uint32) ([]model.StockItem, error)
	ReserveStock(ctx context.Context, sku uint32, stockItem model.StockItem) error
	ListOrder(ctx context.Context, orderID int64) (model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
	ReturnStock(ctx context.Context, userItems model.UserItems) error
	GetOrderStatus(ctx context.Context, orderID int64) (string, error)
	GetOrderItems(ctx context.Context, orderID int64) ([]model.UserItems, error)
	FreeReservedStock(ctx context.Context, userItems model.UserItems) error
	GetOrderUser(ctx context.Context, orderID int64) (int64, error)
}

type Sender interface { // unit test on this
	SendMessage(message model.ProducerMessage) error
}

type Service struct {
	transactionManager TransactionManager
	ordersRepository   OrdersRepository
	sender             Sender
}

func NewService(tx TransactionManager, repo OrdersRepository, sender Sender) *Service {
	return &Service{
		transactionManager: tx,
		ordersRepository:   repo,
		sender:             sender,
	}
}
