package server

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: limit maximum characters
func (server *Server) CreateNews(ctx context.Context, req *pb.CreateNewsRequest) (*pb.NewsIdResponse, error) {
	id, err := server.uc.CreateNews(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error from media/CreateNews")

		if e, ok := err.(utils.BadRequestError); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, err
	}

	return &pb.NewsIdResponse{
		Id: id.String(),
	}, nil
}
