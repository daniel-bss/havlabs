package server

import (
	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/usecases"
	"github.com/daniel-bss/havlabs/internal/media/utils"
)

type Server struct {
	pb.UnimplementedHavlabsMediaServer
	config utils.Config
	uc     usecases.NewsUsecase

	// taskDistributor worker.TaskDistributor
}

func NewGRPCService(config utils.Config, usecase usecases.NewsUsecase) (*Server, error) {
	return &Server{
		config: config,
		uc:     usecase,
	}, nil
}
