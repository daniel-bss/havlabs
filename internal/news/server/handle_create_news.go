package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/google/uuid"
)

func (server *Server) CreateNews(ctx context.Context, req *pb.OneNewsIdRequest) (*pb.NewsIdResponse, error) {
	fmt.Println("create")
	return &pb.NewsIdResponse{
		Id: uuid.New().String(),
	}, nil
}
