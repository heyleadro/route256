package notifications

import (
	"context"
	"fmt"
	"route256/notifications/internal/model"
	"route256/notifications/internal/pkg/logger"
	"time"

	"github.com/opentracing/opentracing-go"
)

// GetUserHistory return user history over a period which is int64 of last hours
func (s *Service) GetUserHistory(ctx context.Context, userID int64, period int64) ([]model.UserNotification, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service/notifications/startconsume")
	defer span.Finish()

	span.SetTag("topic", userID)

	logger.Info("getting history for user %d:", userID)

	key := model.CacheInst{
		UserID: userID,
		Period: period,
	}

	cacheResp, exists := s.cache.Get(key)
	if exists {
		logger.Info("Returning from Cache")
		return cacheResp.([]model.UserNotification), nil
	}

	responseDB, err := s.db.GetUserHistory(userID)
	if err != nil {
		return nil, fmt.Errorf("not exists: %w", err)
	}
	logger.Info("Putting in Cache")

	resultAndCache := findInPeriodBound(responseDB, period)

	s.cache.Add(key, resultAndCache)
	logger.Info("Returning from DB")

	return resultAndCache, nil
}

func findInPeriodBound(in []model.UserNotification, period int64) []model.UserNotification {
	out := make([]model.UserNotification, 0)
	curTime := time.Now()
	for _, val := range in {
		curDiff := curTime.Sub(val.TimeStamp)
		hoursDiff := int64(curDiff.Hours())
		if hoursDiff <= period {
			out = append(out, val)
		}
	}
	return out
}
