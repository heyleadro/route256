package checkout

import (
	"context"
	"errors"
	"route256/checkout/internal/model"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Purchase(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		const (
			userID = int64(1)
		)
		//mockCartRepoResp used for mocking the result returned by GetUserItems method of CartRepo
		mockCartRepoResp := []model.Cart{
			{
				SKU:   uint32(773297411),
				Count: uint16(1),
			},
		}
		// Creating mock interface of CartRepo, which has GetUserItems method called in Service.Purchase
		createCartRepoMock := NewMockCartRepository(t)
		createCartRepoMock.On("GetUserItems", mock.Anything, userID).Return(mockCartRepoResp, nil).Once()

		mockOrderID := int64(1)
		// Creating mock interface of OrderCreator, which has Purchase method called in Service.Purchase
		createOrderMock := NewMockOrderCreator(t)
		createOrderMock.On("Purchase", mock.Anything, userID, mockCartRepoResp).Return(mockOrderID, nil).Once()

		/*
			Act stage, calling Service.Purchase method,
			which inside calls mocked CartRepo.GetUserItems and
			OrderCreator.Purchase
		*/
		clientOrderID, err := (&Service{
			orderCreator: createOrderMock,
			cartRepo:     createCartRepoMock,
		}).Purchase(context.Background(), userID)

		//Assert
		require.NoError(t, err)
		// Service.Purchase returning orderID obtained by OrderCreator.Purchase, so we asserting equality
		require.Equal(t, mockOrderID, clientOrderID)
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
		}).Purchase(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, errMock)
	})

	t.Run("error while purchasing", func(t *testing.T) {
		const (
			userID = int64(1)
		)

		errMock := errors.New("bad resp")

		createCartRepoMock := NewMockCartRepository(t)
		createCartRepoMock.On("GetUserItems", mock.Anything, userID).Return(nil, nil).Once()

		mockOrderID := int64(0)
		createOrderMock := NewMockOrderCreator(t)
		createOrderMock.On("Purchase", mock.Anything, userID, mock.Anything).Return(mockOrderID, errMock).Once()

		// Act
		clientOrderID, err := (&Service{
			orderCreator: createOrderMock,
			cartRepo:     createCartRepoMock,
		}).Purchase(context.Background(), userID)

		// Assert
		require.ErrorIs(t, err, errMock)
		require.Equal(t, clientOrderID, mockOrderID)
	})
}
