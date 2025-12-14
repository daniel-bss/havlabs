package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) CreateUploadURL(ctx context.Context, req *emptypb.Empty) (*pb.SomeResponse, error) {
	fmt.Println("ok")
	return &pb.SomeResponse{Url: "tets"}, nil
}
