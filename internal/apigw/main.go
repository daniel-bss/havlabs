// main.go

package main

import (
	"fmt"
	"net"

	"github.com/daniel-bss/havlabs/internal/apigw/pb"
	"github.com/daniel-bss/havlabs/internal/apigw/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterServiceOneServer(s, &server.Server{})
	reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	defer func(lis net.Listener) {
		if err := lis.Close(); err != nil {
			fmt.Printf("unexpected error: %v\n", err)
		}
	}(lis)

	fmt.Println("okoko")
	s.Serve(lis)
}
