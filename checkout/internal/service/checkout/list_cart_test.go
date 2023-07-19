package checkout

import (
	"context"
	"route256/checkout/internal/model"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ListCart(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		const (
			userID = int64(1)
			Sku    = uint32(773297411)
		)

		//mockCartRepoResp used for mocking the result returned by GetUserItems method of CartRepo
		mockCartRepoResp := []model.Cart{
			{
				SKU:   Sku,
				Count: uint16(1),
			},
		}
		// Creating mock interface of CartRepo, which has GetUserItems method called in Service.ListCart
		createCartRepoMock := NewMockCartRepository(t)
		createCartRepoMock.On("GetUserItems", mock.Anything, userID).Return(mockCartRepoResp, nil).Once()

		mockItem := model.Item{
			Name:  gofakeit.Name(),
			Price: uint32(100),
		}
		// Creating mock interface of CartLister, which has GetProducts method called in Service.ListCart
		listCartMock := NewMockCartLister(t)
		//Must be called once because there is only one item in the test
		listCartMock.On("GetProducts", mock.Anything, Sku).Return(mockItem, nil).Once()
		//Generating ListCart handler like response
		mockListCartDescr := newMockListCartDescr(Sku, mockItem)

		/*
			Act stage, calling Service.ListCart method,
			which inside calls mocked CartRepo.GetUserItems and
			CartLister.GetProducts
		*/
		clientListCartDescr, err := (&Service{
			cartLister: listCartMock,
			cartRepo:   createCartRepoMock,
		}).ListCart(context.Background(), userID)

		//Assert
		require.NoError(t, err)
		// Service.Purchase returning orderID obtained by OrderCreator.Purchase, so we asserting equality
		require.Equal(t, clientListCartDescr, mockListCartDescr)
	})

	t.Run("error while getting user cart from repo", func(t *testing.T) {
		const (
			userID = int64(1)
		)

		errMock := ErrNoUserCart

		createCartRepoMock := NewMockCartRepository(t)
		createCartRepoMock.On("GetUserItems", mock.Anything, userID).Return(nil, errMock).Once()

		// Act
		_, err := (&Service{
			cartRepo: createCartRepoMock,
		}).ListCart(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, errMock)
	})
}

func newMockListCartDescr(Sku uint32, item model.Item) *model.ListCartDescr {
	mockListCartDescr := &model.ListCartDescr{
		ItemsDescr: []model.ItemDescr{
			{
				CartDescr: model.Cart{
					SKU:   Sku,
					Count: uint16(1),
				},
				UserItemDescr: item,
			},
		},
	}
	totalPrice := uint32(0)
	for _, val := range mockListCartDescr.ItemsDescr {
		totalPrice += val.UserItemDescr.Price
	}
	mockListCartDescr.TotalPrice = totalPrice
	return mockListCartDescr
}
