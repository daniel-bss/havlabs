package server

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*pb.ConfirmUploadResponse, error) {
	uploadConfirmation, err := server.uc.ConfirmUpload(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("error from media/ConfirmUpload")

		if e, ok := err.(utils.BadRequestError); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, err
	}

	return &pb.ConfirmUploadResponse{
		MediaId: uploadConfirmation.MediaId,
		Status:  uploadConfirmation.Status,
	}, nil
}
