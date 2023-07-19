package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "ozon",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   []float64{0.5, 0.9, 0.99},
	},
		[]string{
			"status",
			"method",
		},
	)

	RequestsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
		},
		[]string{"status", "method"},
	)
)
