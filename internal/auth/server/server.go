package server

import (
	"fmt"

	"github.com/daniel-bss/havlabs/auth/pb"
	"github.com/daniel-bss/havlabs/auth/token"
	"github.com/daniel-bss/havlabs/auth/utils"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedHavlabsAuthServer
	config     utils.Config
	tokenMaker token.Maker
	// store           db.Store
	// taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func New(config utils.Config, store any, taskDistributor any) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		// store:           store,
		// taskDistributor: taskDistributor,
	}

	return server, nil
}
