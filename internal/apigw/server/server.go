package server

import (
	"context"

	"github.com/daniel-bss/havlabs/internal/apigw/pb"
)

type Server struct {
	pb.UnimplementedServiceOneServer
}

func (s *Server) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Msg: "Hello, " + req.Name,
	}, nil
}
