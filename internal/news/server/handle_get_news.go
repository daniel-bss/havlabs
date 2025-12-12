package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) GetAllNews(ctx context.Context, req *emptypb.Empty) (*pb.GetAllNewsResponse, error) {
	// mtdt := server.extractMetadata(ctx)
	// if md, ok := metadata.FromIncomingContext(ctx); ok {
	// 	fmt.Println(md)
	// }

	fmt.Println("get all")

	// return &pb.NewsResponse{Title: "TEST"}, nil
	return &pb.GetAllNewsResponse{
		News: []*pb.GetOneNewsResponse{
			{
				Title: "hehe1",
			},
			{
				Title: "hehe2",
			},
		},
	}, nil
}

func (server *Server) GetOneNews(ctx context.Context, req *pb.OneNewsIdRequest) (*pb.GetOneNewsResponse, error) {
	fmt.Println("get one")
	return &pb.GetOneNewsResponse{Title: "TEST"}, nil
}
