package product

import (
	"log"
	"route256/checkout/pkg/product_v1"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	RPS   = 10
	Burst = 10
)

type Client struct {
	product_v1.ProductServiceClient
	*rate.Limiter
}

func NewClient(port string) *Client {
	ratelimiter := rate.NewLimiter(rate.Limit(RPS), Burst)

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	productServiceClient := product_v1.NewProductServiceClient(conn)

	return &Client{
		ProductServiceClient: productServiceClient,
		Limiter:              ratelimiter,
	}
}
