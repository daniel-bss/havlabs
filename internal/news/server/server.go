package server

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	pb.UnimplementedHavlabsNewsServer
	config utils.Config
	// store      db.Store
	// taskDistributor worker.TaskDistributor
}

func NewGRPCService(config utils.Config) (*Server, error) {
	server := &Server{
		config: config,
		// store:      store,
	}

	return server, nil
}

func (server *Server) GetNews(ctx context.Context, req *pb.NewsRequest) (*pb.NewsResponse, error) {
	// mtdt := server.extractMetadata(ctx)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Println("HEHE")
		fmt.Println(md)
	}
	fmt.Println()
	return &pb.NewsResponse{Title: "HEHE"}, nil
}
