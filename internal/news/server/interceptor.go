package server

import (
	"context"

	"google.golang.org/grpc"
)

func CustomInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		return handler(ctx, req)
	}
}
