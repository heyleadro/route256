package logging

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("method: %v, req: %v\n", info.FullMethod, req)
	resp, err = handler(ctx, req)
	log.Printf("resp: %v, err: %v\n", resp, err)
	return resp, err
}
