package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/google/uuid"
)

func (server *Server) UpdateNews(ctx context.Context, req *pb.OneNewsIdRequest) (*pb.NewsIdResponse, error) {
	fmt.Println("update")
	return &pb.NewsIdResponse{
		Id: uuid.New().String(),
	}, nil
}
