package model

import (
	"fmt"
	"sync"
	"time"
)

// instead of postgres for example, !it is not cache!
type UserNotificationHistory struct {
	His map[int64][]UserNotification
	mu  sync.RWMutex
}

type UserNotification struct {
	OrderID   int64
	TimeStamp time.Time
}

type CacheInst struct {
	UserID int64
	Period int64
}

type Info struct {
	UserID    int64
	OrderID   int64
	TimeStamp time.Time
}

func (m *UserNotificationHistory) AddToDB(userID int64, orderID int64, timeStamp time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	userNotification, exists := m.His[userID]
	if !exists {
		userNotification = make([]UserNotification, 0)
	}

	userNotification = append(userNotification, UserNotification{
		OrderID:   orderID,
		TimeStamp: timeStamp,
	})

	m.His[userID] = userNotification

	return nil
}

func (m *UserNotificationHistory) GetUserHistory(userID int64) ([]UserNotification, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	userNotification, exists := m.His[userID]
	if !exists {
		return nil, fmt.Errorf("no user")
	}

	return userNotification, nil
}

func NewModel() *UserNotificationHistory {
	return &UserNotificationHistory{
		His: make(map[int64][]UserNotification),
	}
}
