package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"google.golang.org/grpc/metadata"
)

func (server *Server) GetNews(ctx context.Context, req *pb.NewsRequest) (*pb.NewsResponse, error) {
	// mtdt := server.extractMetadata(ctx)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Println(md)
	}

	fmt.Println()

	return &pb.NewsResponse{Title: "TEST"}, nil
}
