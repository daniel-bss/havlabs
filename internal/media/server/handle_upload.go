package server

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) CreateUploadURL(ctx context.Context, req *emptypb.Empty) (*pb.SomeResponse, error) {
	uploadUrl, err := server.uc.CreateUpload(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.SomeResponse{Url: uploadUrl}, nil
}
