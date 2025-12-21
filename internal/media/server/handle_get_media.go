package server

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GET ID
func (server *Server) GetMediaById(ctx context.Context, req *pb.GetOneMediaByIdRequest) (*pb.OneMediaResponse, error) {
	media, err := server.uc.GetMediaById(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error from media/GetMediaById")

		if e, ok := err.(utils.BadRequestError); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, err
	}

	return &pb.OneMediaResponse{
		Id: media.Id.String(),
	}, nil
}

// GET URL STRING
func (server *Server) GetMediaURL(ctx context.Context, req *pb.GetOneMediaByIdRequest) (*pb.OneMediaURLResponse, error) {
	mediaURL, err := server.uc.GetMediaURLById(ctx, utils.ParseInt(server.config.PresignedGETUrlDurationMinutes), req)
	if err != nil {
		log.Error().Err(err).Msg("error from media/GetMediaURLById")

		if e, ok := err.(utils.BadRequestError); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, err
	}

	return &pb.OneMediaURLResponse{
		Url: mediaURL,
	}, nil
}
