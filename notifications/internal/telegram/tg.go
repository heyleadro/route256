package telegram

import (
	"fmt"
	"route256/notifications/internal/pkg/logger"
	"time"

	tele "gopkg.in/telebot.v3"
)

type NotifyBot struct {
	bot     tele.Bot
	channel tele.Chat
}

func NewNotifyBot(token string, chatId int64) (*NotifyBot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	logger.Info("created bot\n")

	go bot.Start()

	logger.Info("started bot\n")

	channel := tele.Chat{
		ID:   chatId,
		Type: tele.ChatChannel,
	}

	return &NotifyBot{
		bot:     *bot,
		channel: channel,
	}, nil
}

func (b *NotifyBot) SendMSG(text string) error {
	logger.Info("Sending message")

	_, err := b.bot.Send(&b.channel, text)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	logger.Info("sent message")

	return nil
}
