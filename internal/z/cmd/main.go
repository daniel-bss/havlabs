package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/daniel-bss/havlabs/internal/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect to gRPC server")
	}
	defer conn.Close()

	client := pb.NewHavlabsAuthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rsp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "a",
		Password: "x",
	})
	if err != nil {
		fmt.Printf("failed to Login: %s\n", err.Error())
		s := status.Convert(err)
		fmt.Println(">>", s, s.Code(), s.Details())
		for _, d := range s.Details() {
			fmt.Println(d)
			fmt.Println(reflect.TypeOf(d))
			fmt.Println()
			// switch info := d.(type) {
			// case :
			// 	log.Printf("Quota failure: %s", info)
			// default:
			// 	log.Printf("Unexpected type: %s", info)
			// }
		}
	}

	fmt.Println("Received:")
	fmt.Println(rsp)

}
