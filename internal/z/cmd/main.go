package main

import (
	"fmt"

	_ "github.com/daniel-bss/havlabs/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect to gRPC server")
	}
	defer conn.Close()

}
