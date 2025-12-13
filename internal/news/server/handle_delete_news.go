package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/google/uuid"
)

func (server *Server) DeleteNews(ctx context.Context, req *pb.GetOneNewsByIdRequest) (*pb.NewsIdResponse, error) {
	fmt.Print("delete")
	return &pb.NewsIdResponse{
		Id: uuid.New().String(),
	}, nil
}
