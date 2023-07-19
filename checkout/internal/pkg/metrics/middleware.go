package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()

	h, err := handler(ctx, req)

	status := status.Code(err)

	HistogramResponseTime.WithLabelValues(status.String(), info.FullMethod).Observe(time.Since(start).Seconds())
	RequestsCounter.WithLabelValues(status.String(), info.FullMethod).Inc()

	return h, err
}
