package checkout

import (
	"context"
	"fmt"
	"route256/checkout/internal/model"
	"route256/checkout/internal/pkg/logger"
	"route256/checkout/internal/pkg/tracer"
	"route256/checkout/internal/pkg/workerpool"
	"sync"

	"github.com/opentracing/opentracing-go"
)

// MaxWorkers limits at once requests
const MaxWorkers = 5

// struct to store ListCartDescr result and err of each request
// used for channel inside wp.Exec func
type ItemDescrWithErr struct {
	itemDescr model.ItemDescr
	err       error
}

/*
ListCart list get_product response from productservice
*/
func (s *Service) ListCart(ctx context.Context, user int64) (*model.ListCartDescr, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/checkout/listcart")
	defer span.Finish()

	span.SetTag("user_id", user)

	// getting user items in cart from db
	cartItems, err := s.cartRepo.GetUserItems(ctx, user)
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, ErrNoUserCart)
	}

	wp := workerpool.NewWP(MaxWorkers) // initializing worker pool

	itemsChan := make(chan ItemDescrWithErr, len(cartItems)) // initializing channel to store results from wp

	wg := sync.WaitGroup{} // wg to wait until all cartItems processed

	for _, v := range cartItems {
		wg.Add(1)

		curItem := v // namespace

		err := wp.Exec(ctx, func(ctx context.Context) { // implementing func which will be called in goroutine inside exec
			defer wg.Done()
			item, err := s.cartLister.GetProducts(ctx, v.SKU) // calling getproducts
			itemsChan <- ItemDescrWithErr{                    // writing into channel
				itemDescr: model.ItemDescr{
					CartDescr:     curItem,
					UserItemDescr: item,
				},
				err: err, // storing err, to be processed later
			}
		})
		if err != nil {
			return nil, tracer.MarkSpanWithError(ctx, fmt.Errorf("wp exec: %w", err))
		}
	}
	wg.Wait()

	items := make([]model.ItemDescr, 0, len(cartItems))
	var totalPrice uint32

	for i := 0; i < len(cartItems); i++ { // traversing through channel
		curItem := <-itemsChan  // reading from result channel
		if curItem.err != nil { // checking each response of getproduct on errors
			return nil, tracer.MarkSpanWithError(ctx, fmt.Errorf("get product response: %w", curItem.err))
		}

		// generating final result
		items = append(items, curItem.itemDescr)
		totalPrice += curItem.itemDescr.UserItemDescr.Price * (uint32(curItem.itemDescr.CartDescr.Count))
	}

	// some logs
	logger.Info("User's %d items", user)
	for _, item := range items {
		logger.Info("Item: %s, Price: %d", item.UserItemDescr.Name, item.UserItemDescr.Price)
	}
	logger.Info("Total price: %d", totalPrice)

	return &model.ListCartDescr{
		TotalPrice: totalPrice,
		ItemsDescr: items,
	}, nil
}
