package notifications

import (
	"route256/notifications/internal/service/notifications"
	"route256/notifications/pkg/notifications_v1"
)

type Service struct {
	notifications_v1.UnimplementedNotificationsServer
	impl *notifications.Service
}

func NewNotificationsServer(impl *notifications.Service) notifications_v1.NotificationsServer {
	return &Service{impl: impl}
}
