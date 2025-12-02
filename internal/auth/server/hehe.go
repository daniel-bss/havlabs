package server

import (
	"context"
	"fmt"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcHehe(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	// startTime := time.Now()
	result, err := handler(ctx, req)
	// duration := time.Since(startTime)

	// statusCode := codes.Unknown
	// if st, ok := status.FromError(err); ok {
	// 	statusCode = st.Code()
	// }

	fmt.Println("OKCAKKK")
	fmt.Println(result)
	fmt.Println(reflect.TypeOf(result))

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	fmt.Print(statusCode)

	return result, err
}
