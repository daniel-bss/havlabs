package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/media/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/media/dtos"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

func (uc *newsUsecaseImpl) ConfirmUpload(ctx context.Context, req *pb.ConfirmUploadRequest) (*dtos.ConfirmUpload, error) {
	// validate UUID
	id, err := uuid.Parse(req.MediaId)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid media id")
	}
	media, err := uc.store.GetMediaById(ctx, id)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid media id")
	}

	// handle if media has been promoted from staging
	stagingBucket := uc.minioManager.GetStagingBucketName()
	if media.Bucket != stagingBucket {
		return nil, utils.NewBadRequestError("invalid media id")
	}

	uploadedStatusArg := db.UpdateUploadStatusParams{
		ID:     media.ID,
		Status: db.StatusTypesUploaded,
	}
	failedStatusArg := db.UpdateUploadStatusParams{
		ID:     media.ID,
		Status: db.StatusTypesFailed,
	}

	// update status
	_, err = uc.store.UpdateUploadStatus(ctx, uploadedStatusArg)
	if err != nil {
		return nil, err
	}

	// prepare client, bucket, staging bucket, objectKey
	client := uc.minioManager.GetClient()
	bucket := uc.minioManager.GetBucketName()
	objectKey := fmt.Sprintf("%s/%s", media.Purpose, media.ID)
	opts := minio.StatObjectOptions{}

	// get object from staging
	reader, err := client.GetObject(ctx, stagingBucket, objectKey, opts)
	if err != nil {
		log.Error().Err(err).Msg("error from GetObject 1")
		return nil, err
	}
	defer reader.Close()

	// MIME sniffing
	mime, err := utils.DetectMIME(reader)
	if err != nil {
		log.Error().Err(err).Msg("error when MIME Sniffing")
		return nil, err
	}
	if !utils.IsValidContentType(mime) {
		_, err = uc.store.UpdateUploadStatus(ctx, failedStatusArg)
		if err != nil {
			return nil, err
		}
		return nil, utils.NewBadRequestError(fmt.Sprintf("invalid mime: %s", mime))
	}

	// checksum
	checksum := utils.GetChecksum(reader)

	// metadata
	// 1. new reader, fresh
	reader, err = client.GetObject(ctx, stagingBucket, objectKey, opts)
	if err != nil {
		log.Error().Err(err).Msg("error from GetObject 2")
		return nil, err
	}
	defer reader.Close()

	// 2. get the image metadata
	imgConfig, fileFormat, err := image.DecodeConfig(reader)
	if err != nil {
		log.Error().Err(err).Msg("error from image.DecodeConfig")
		return nil, err
	}
	imgMetadata, err := json.Marshal(map[string]any{
		"width":  imgConfig.Width,
		"height": imgConfig.Height,
		"format": fileFormat,
	})
	if err != nil {
		log.Error().Err(err).Msg("error when marshalling metadata")
		return nil, err
	}

	// bytes size
	obj, err := client.StatObject(ctx, stagingBucket, objectKey, opts)
	if err != nil {
		log.Error().Err(err).Msg("error from StatObject")
		return nil, err
	}
	if obj.Size > utils.MaxImageSize {
		_, err = uc.store.UpdateUploadStatus(ctx, failedStatusArg)
		if err != nil {
			return nil, err
		}
		return nil, utils.NewBadRequestError(fmt.Sprintf("image size too large. maximum: %d B", utils.MaxImageSize))
	}

	// MARK: end of validation & preprocessing
	// begin promoting bucket

	srcObj := minio.CopySrcOptions{
		Bucket: stagingBucket,
		Object: objectKey,
	}

	dstObj := minio.CopyDestOptions{
		Bucket: bucket,
		Object: fmt.Sprintf("%s.%s", objectKey, fileFormat), // attach actual file format
	}

	info, err := client.CopyObject(ctx, dstObj, srcObj)
	if err != nil {
		log.Error().Err(err).Msg("error when moving object")
		return nil, err
	}

	err = client.RemoveObject(
		ctx,
		stagingBucket,
		objectKey,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		log.Error().Err(err).Msg("error when RemoveObject")
		return nil, err
	}

	// update media db
	arg := db.ConfirmGoodUploadParams{
		ID:                media.ID,
		Bucket:            bucket,
		ActualContentType: pgtype.Text{String: mime, Valid: true},
		SizeBytes:         pgtype.Int8{Int64: obj.Size, Valid: true},
		Checksum:          pgtype.Text{String: checksum, Valid: true},
		Metadata:          imgMetadata,
	}
	mediaId, err := uc.store.ConfirmGoodUpload(ctx, arg)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("uploaded media of id %s, etag %s", mediaId.String(), info.ETag)

	return &dtos.ConfirmUpload{
		MediaId: mediaId.String(),
		Status:  pb.StatusEnum_ready,
	}, nil
}
