package notifications

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_StartConsume(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		const (
			topic = "orders"
		)
		botMock := NewMockBot(t)
		botMock.On("SendMSG", mock.Anything).Return(nil).Maybe()

		consumerMock := NewMockConsumer(t)
		consumerMock.On("Subscribe", topic, mock.AnythingOfType("func(string) error")).Return(nil).Once()

		serviceErr := (&Service{
			tgBot: botMock,
		}).StartConsume(topic)

		require.NoError(t, serviceErr)
	})

	t.Run("error", func(t *testing.T) {
		const (
			topic = "orders"
		)
		errMock := errors.New("err some")

		botMock := NewMockBot(t)
		botMock.On("SendMSG", mock.Anything).Return(nil).Maybe()

		consumerMock := NewMockConsumer(t)
		consumerMock.On("Subscribe", topic, mock.AnythingOfType("func(string) error"), mock.Anything).Return(errMock).Once()

		serviceErr := (&Service{
			tgBot: botMock,
		}).StartConsume(topic)

		require.ErrorIs(t, serviceErr, errMock)
	})
}
