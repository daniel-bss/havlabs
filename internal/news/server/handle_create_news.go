package server

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
)

// TODO: limit maximum characters
func (server *Server) CreateNews(ctx context.Context, req *pb.CreateNewsRequest) (*pb.NewsIdResponse, error) {
	id, err := server.uc.CreateNews(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.NewsIdResponse{
		Id: id.String(),
	}, nil
}
