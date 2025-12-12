package server

import (
	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/usecases"
	"github.com/daniel-bss/havlabs/internal/news/utils"
)

type Server struct {
	pb.UnimplementedHavlabsNewsServer
	config utils.Config
	uc     usecases.NewsUsecase

	// taskDistributor worker.TaskDistributor
}

func NewGRPCService(config utils.Config, uc usecases.NewsUsecase) (*Server, error) {
	return &Server{
		config: config,
		uc:     uc,
	}, nil
}
