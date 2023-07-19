package checkout

import (
	"context"
	"errors"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/product"
	"route256/checkout/internal/config"
	"route256/checkout/internal/model"
)

type StockChecker interface {
	Stocks(ctx context.Context, sku uint32) ([]model.Stock, error)
}

type OrderCreator interface {
	Purchase(ctx context.Context, user int64, items []model.Cart) (int64, error)
}

type CartLister interface {
	GetProducts(ctx context.Context, sku uint32) (model.Item, error)
}

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type CartRepository interface {
	AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error
	GetUserSkuItemCount(ctx context.Context, user int64, sku uint32) (uint16, error)
	GetUserItems(ctx context.Context, user int64) ([]model.Cart, error)
}

type Service struct {
	stockChecker       StockChecker
	orderCreator       OrderCreator
	cartRepo           CartRepository
	transactionManager TransactionManager
	cartLister         CartLister
}

func NewService(tx TransactionManager, repo CartRepository) *Service {
	return &Service{
		stockChecker:       loms.NewClient(config.AppConfig.Services.Loms),
		orderCreator:       loms.NewClient(config.AppConfig.Services.Loms),
		cartLister:         product.NewClient(config.AppConfig.Services.Product),
		transactionManager: tx,
		cartRepo:           repo,
	}
}

var (
	ErrStockInsufficient = errors.New("stock insufficient")
	ErrCartInsufficient  = errors.New("cart insufficient")
	ErrNoCarts           = errors.New("no carts")
	ErrNoUserCart        = errors.New("no cart for user")
)
