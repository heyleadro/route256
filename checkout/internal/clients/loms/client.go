package loms

import (
	"log"
	"route256/checkout/pkg/loms_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	loms_v1.LomsClient
}

func NewClient(port string) *Client {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	lomsClient := loms_v1.NewLomsClient(conn)

	return &Client{LomsClient: lomsClient}
}
