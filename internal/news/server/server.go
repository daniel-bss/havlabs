package server

import (
	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/utils"
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
