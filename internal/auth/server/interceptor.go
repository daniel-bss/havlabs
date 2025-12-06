package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func Interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	// startTime := time.Now()
	// _, err := handler(ctx, req)
	// duration := time.Since(startTime)

	// statusCode := codes.Unknown
	// if st, ok := status.FromError(err); ok {
	// 	statusCode = st.Code()
	// }

	// fmt.Println("OKCAKKK")
	// fmt.Println(result)
	// fmt.Println(reflect.TypeOf(result))

	// statusCode := codes.Unknown
	// if st, ok := status.FromError(err); ok {
	// 	statusCode = st.Code()
	// }
	// fmt.Print(statusCode)

	fmt.Println("TODO: interceptor")

	// return nil, status.Errorf(codes.Internal, "wah unsupported message type: ")

	return handler(ctx, req)

}
