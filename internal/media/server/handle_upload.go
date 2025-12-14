package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) CreateUploadURL(ctx context.Context, req *emptypb.Empty) (*pb.SomeResponse, error) {
	fmt.Println("ok")

	uploadUrl := server.uc.CreateUpload(ctx)

	return &pb.SomeResponse{Url: uploadUrl}, nil
}
