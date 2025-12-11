package server

import (
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/auth/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/auth/token"
	"github.com/daniel-bss/havlabs/internal/auth/utils"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedHavlabsAuthServer
	config     utils.Config
	tokenMaker *token.JWTMaker
	store      db.Store
	// taskDistributor worker.TaskDistributor
}

// TODO: github.com/theckman/go-pwnedpasswords
// TODO: github.com/theifedayo/go-dumb-password

// NewServer creates a new gRPC server.
func NewGRPCService(config utils.Config, store db.Store, taskDistributor any) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker()
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
		store:      store,
		// taskDistributor: taskDistributor,
	}

	return server, nil
}
