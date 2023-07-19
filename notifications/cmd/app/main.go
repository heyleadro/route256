package main

import (
	"context"
	"net"
	api "route256/notifications/internal/api/notifications"
	"route256/notifications/internal/config"
	"route256/notifications/internal/kafka"
	"route256/notifications/internal/pkg/logger"
	"route256/notifications/internal/pkg/tracer"
	"route256/notifications/internal/service/notifications"
	"route256/notifications/internal/telegram"
	"route256/notifications/pkg/notifications_v1"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = ":50053"

func main() {
	if err := tracer.InitGlobal(notifications.ServiceName); err != nil {
		logger.Fatal("ERR: ", err)
	}

	err := config.Init()
	if err != nil {
		logger.Fatal("ERR: ", err)
	}

	kafkaConsumer, err := kafka.NewConsumer(config.AppConfig.Kafka.Brokers)
	if err != nil {
		logger.Fatal("kafka: ", err)
	}

	tgBot, err := telegram.NewNotifyBot(
		config.AppConfig.Bot.Token,
		config.AppConfig.Bot.ChatID,
	)
	if err != nil {
		logger.Fatal("tg bot: ", err)
	}

	notifications := notifications.NewService(kafkaConsumer, tgBot)

	go func() {
		err = notifications.StartConsume(config.AppConfig.Kafka.Topic)
		if err != nil {
			logger.Fatal("start consume: ", err)
		}
	}()
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			logger.Fatal("no notify ", err)
			return err
		}

		s := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				logger.MiddlewareGRPC,
				tracer.MiddlewareGRPC,
			),
		)

		reflection.Register(s)

		notifications_v1.RegisterNotificationsServer(s, api.NewNotificationsServer(notifications))

		logger.Info("server listening at %v", lis.Addr())

		if err = s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
		return nil
	})

	go func() {
		err = g.Wait()
		logger.Info("server stopped listening at 50053")

		if err != nil {
			logger.Fatal("g wait")
		}
	}()
	<-context.TODO().Done()
}
