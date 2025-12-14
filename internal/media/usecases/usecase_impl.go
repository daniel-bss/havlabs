package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

type newsUsecaseImpl struct {
	minioClient *minio.Client
}

func New(mc *minio.Client) NewsUsecase {
	return &newsUsecaseImpl{
		minioClient: mc,
	}
}

func (uc *newsUsecaseImpl) CreateUpload(ctx context.Context) string {
	bucketName := "mybucket"
	// location := ""

	// err := uc.minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	// if err != nil {
	// 	exists, errBucketExists := uc.minioClient.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		log.Error().Err(errors.New("bucket already exists")).Msgf("bucket %s", bucketName)
	// 	} else {
	// 		log.Fatal().Err(err).Msg("error when making bucket")
	// 	}
	// }

	log.Info().Msgf("successfully created bucket: %s", bucketName)

	objectKey := fmt.Sprintf("images/%s.png", uuid.NewString())
	url, err := uc.minioClient.PresignedPutObject(ctx, bucketName, objectKey, time.Second*90)
	if err != nil {
		log.Error().Err(err).Msg("failed to PresignedPutObject")
	}

	return url.String()
}
