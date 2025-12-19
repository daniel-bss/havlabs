package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/media/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/media/dtos"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *newsUsecaseImpl) CreateUploadSession(ctx context.Context, durationMinutes int, req *pb.CreateUploadSessionRequest) (*dtos.UploadSession, error) {
	if !utils.IsValidPurpose(req.Purpose) {
		return nil, status.Error(codes.InvalidArgument, "invalid purpose")
	}

	if !utils.IsValidContentType(req.ContentType) {
		return nil, status.Error(codes.InvalidArgument, "invalid content type")
	}

	mediaId := uuid.New()
	stagingBucket := uc.minioManager.GetStagingBucketName() // staging
	objectKey := fmt.Sprintf("%s/%s", req.Purpose, mediaId.String())
	client := uc.minioManager.GetClient()

	uploadUrl, err := client.PresignedPutObject(ctx, stagingBucket, objectKey, time.Minute*time.Duration(durationMinutes))
	if err != nil {
		log.Error().Err(err).Msg("failed when PresignedPutObject")
		return nil, err
	}

	// get metadata
	username, err := utils.ExtractMetadataFromContextWithKey("x-username", ctx)
	if err != nil {
		return nil, err
	}

	// INSERT to db
	arg := db.CreateUploadParams{
		ID:                  mediaId,
		OwnerUsername:       username,
		FileName:            req.FileName,
		Purpose:             req.Purpose,
		DeclaredContentType: req.ContentType,
		Bucket:              stagingBucket,
		ObjectKey:           objectKey,
	}
	mediaId, err = uc.store.CreateUpload(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &dtos.UploadSession{
		MediaId:   mediaId.String(),
		UploadUrl: uploadUrl.String(),
	}, nil
}
