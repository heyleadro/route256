package loms

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/model"
)

// add unit test
func (s *Service) Notify(ctx context.Context, userID int64, orderID int64, status string) error {
	log.Printf("notifying about order: %d", orderID)

	err := s.sender.SendMessage(model.ProducerMessage{
		UserID:  userID,
		OrderID: orderID,
		Status:  status,
	})
	if err != nil {
		return fmt.Errorf("notify : %w", err)
	}

	return nil
}
