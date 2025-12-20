package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) GetAllNews(ctx context.Context, req *emptypb.Empty) (*pb.ListNewsResponse, error) {
	fmt.Println("get all")

	return &pb.ListNewsResponse{
		News: []*pb.OneNewsResponse{
			{
				Title: "hehe1",
			},
			{
				Title: "hehe2",
			},
		},
	}, nil
}

func (server *Server) GetOneNews(ctx context.Context, req *pb.GetOneNewsByIdRequest) (*pb.OneNewsResponse, error) {
	fmt.Println("get one")
	return &pb.OneNewsResponse{Title: "TEST"}, nil
}
