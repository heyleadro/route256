package notifications

import (
	"route256/libs/cache"
	"route256/notifications/internal/model"
	"time"
)

const ServiceName = "notifications"

type Consumer interface {
	Subscribe(topic string, handle func(text string) error, ch chan model.Info) error
}

type Bot interface {
	SendMSG(text string) error
}

type DB interface {
	AddToDB(userID int64, orderID int64, timeStamp time.Time) error
	GetUserHistory(userID int64) ([]model.UserNotification, error)
}

type Service struct {
	consumer Consumer
	tgBot    Bot
	db       DB
	cache    *cache.Cache
}

const cap = 200

func NewService(c Consumer, bot Bot) *Service {
	return &Service{
		consumer: c,
		tgBot:    bot,
		db:       model.NewModel(),
		cache:    cache.NewCache(cap),
	}
}
