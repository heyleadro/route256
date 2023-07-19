package loms

import (
	"context"
	"errors"
	"route256/loms/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Notify(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var (
			msg = model.ProducerMessage{
				UserID:  1,
				OrderID: 1,
				Status:  "OK",
			}
		)

		sendMock := NewMockSender(t)
		sendMock.On("SendMessage", msg).Return(nil).Once()

		serviceErr := (&Service{
			sender: sendMock,
		}).Notify(context.Background(), msg.UserID, msg.OrderID, msg.Status)

		require.NoError(t, serviceErr)
	})

	t.Run("error", func(t *testing.T) {
		var (
			msg = model.ProducerMessage{
				OrderID: 1,
				Status:  "OK",
			}
			errMock = errors.New("some err")
		)

		sendMock := NewMockSender(t)
		sendMock.On("SendMessage", msg).Return(errMock).Once()

		serviceErr := (&Service{
			sender: sendMock,
		}).Notify(context.Background(), msg.UserID, msg.OrderID, msg.Status)

		require.ErrorIs(t, serviceErr, errMock)
	})
}
