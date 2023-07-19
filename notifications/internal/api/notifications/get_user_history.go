package notifications

import (
	"context"
	"route256/notifications/internal/model"
	"route256/notifications/internal/pkg/logger"
	"route256/notifications/pkg/notifications_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetUserHistory(ctx context.Context, req *notifications_v1.GetUserHistoryRequest) (*notifications_v1.GetUserHistoryResponse, error) {
	logger.Info("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	userNotification, err := s.impl.GetUserHistory(ctx, req.User, int64(req.Period))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return bindToResp(userNotification), nil
}

func bindToResp(in []model.UserNotification) *notifications_v1.GetUserHistoryResponse {
	items := make([]*notifications_v1.HisItem, 0, len(in))

	for _, val := range in {
		items = append(items, &notifications_v1.HisItem{
			Order: val.OrderID,
			Time:  val.TimeStamp.Format("2006-01-02 15:04:05"),
		})
	}

	return &notifications_v1.GetUserHistoryResponse{
		His: items,
	}
}
