package notifications

import (
	"context"
	"route256/notifications/internal/model"
	"route256/notifications/internal/pkg/logger"

	"github.com/opentracing/opentracing-go"
)

func (s *Service) StartConsume(topic string) error {
	span, _ := opentracing.StartSpanFromContext(context.Background(), "service/notifications/startconsume")
	defer span.Finish()

	span.SetTag("topic", topic)

	logger.Info("start consuming topic %s:", topic)

	ch := make(chan model.Info)
	go func() {
		err := s.consumer.Subscribe(topic, s.tgBot.SendMSG, ch)
		if err != nil {
			logger.Fatalf("subscribe: ", err)
		}
	}()
	for {
		select {
		case resp := <-ch:
			err := s.db.AddToDB(resp.UserID, resp.OrderID, resp.TimeStamp)
			response, err := s.db.GetUserHistory(resp.UserID)
			for _, val := range response {
				logger.Info(val.OrderID, val.TimeStamp)
			}
			if err != nil {
				break
			}
		default:
			continue
		}
	}

	return nil
}
