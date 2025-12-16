package client

import (
	"context"
	"fmt"

	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

type MinioManager interface {
	GetBucketName() string
	GetClient() *minio.Client
}

type minioManagerImpl struct {
	client     *minio.Client
	bucketName string
}

func NewMinioManager(config utils.Config) MinioManager {
	endpoint := fmt.Sprintf("%s:%s", config.MinioHost, config.MinioServerPort)
	accessKeyID := config.MinioRootUser
	secretAccessKey := config.MinioRootPassword

	useSSL := utils.ParseBool(config.MinioUseSSL)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create MinIO client")
	}

	log.Info().Msg("successfully init MinIO client")

	bucketName := config.MinioBucketName
	location := ""

	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Info().Err(err).Msgf("bucket %s already exists", bucketName)
		} else {
			log.Fatal().Err(err).Msg("error when making bucket")
		}
	}

	log.Info().Msgf("successfully created bucket: %s", bucketName)

	return &minioManagerImpl{
		client:     minioClient,
		bucketName: bucketName,
	}
}

func (m *minioManagerImpl) GetBucketName() string {
	return m.bucketName
}

func (m *minioManagerImpl) GetClient() *minio.Client {
	return m.client
}
