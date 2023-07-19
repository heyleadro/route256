package notifications

import (
	"context"
	"route256/libs/cache"
	"route256/notifications/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GetUserHistory(t *testing.T) {

	t.Run("success from cache", func(t *testing.T) {
		mockKey := model.CacheInst{
			UserID: 1,
			Period: 1,
		}

		mockHistory := []model.UserNotification{
			{
				OrderID:   1,
				TimeStamp: time.Now(),
			},
		}
		mockCache := cache.NewCache(2)
		mockCache.Add(mockKey, mockHistory)
		// consumerMock := NewMockDB(t)
		// consumerMock.On("GetUserHistory", userID).Return(mockHistory, nil)

		serviceResp, serviceErr := (&Service{
			cache: mockCache,
		}).GetUserHistory(context.Background(), mockKey.UserID, mockKey.Period)

		require.NoError(t, serviceErr)

		mockResp, exists := mockCache.Get(mockKey)
		require.Equal(t, exists, true)

		require.ElementsMatch(t, serviceResp, mockResp.([]model.UserNotification))
	})

	t.Run("success from db and then from cache", func(t *testing.T) {
		mockKey := model.CacheInst{
			UserID: 1,
			Period: 1,
		}

		mockHistory := []model.UserNotification{
			{
				OrderID:   1,
				TimeStamp: time.Now(),
			},
		}
		mockCache := cache.NewCache(2)
		// mockCache.Add(mockKey, mockHistory)
		dbMock := NewMockDB(t)
		dbMock.On("GetUserHistory", mockKey.UserID).Return(mockHistory, nil).Once()

		service := &Service{
			db:    dbMock,
			cache: mockCache,
		}
		serviceResp, serviceErr := service.GetUserHistory(context.Background(), mockKey.UserID, mockKey.Period)

		require.NoError(t, serviceErr)

		mockResp, exists := mockCache.Get(mockKey)
		require.Equal(t, exists, true)

		require.ElementsMatch(t, serviceResp, mockResp.([]model.UserNotification))

		serviceResp, serviceErr = service.GetUserHistory(context.Background(), mockKey.UserID, mockKey.Period)

		require.NoError(t, serviceErr)

		require.Equal(t, exists, true)

		require.ElementsMatch(t, serviceResp, mockResp.([]model.UserNotification))
	})
}
