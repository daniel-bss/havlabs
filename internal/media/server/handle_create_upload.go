package server

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUploadSession(ctx context.Context, req *pb.CreateUploadSessionRequest) (*pb.CreateUploadSessionResponse, error) {
	uploadSession, err := server.uc.CreateUploadSession(ctx, utils.ParseInt(server.config.PresignedPUTUrlDurationMinutes), req)
	if err != nil {
		log.Error().Err(err).Msg("error from media/CreateUploadSession")

		if e, ok := err.(utils.BadRequestError); ok {
			return nil, status.Error(codes.InvalidArgument, e.Error())
		}
		return nil, err
	}

	return &pb.CreateUploadSessionResponse{
		MediaId:   uploadSession.MediaId,
		UploadUrl: uploadSession.UploadUrl,
	}, nil
}
