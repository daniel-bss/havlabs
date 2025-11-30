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

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	fmt.Println("okoko")
	s.Serve(lis)
}
