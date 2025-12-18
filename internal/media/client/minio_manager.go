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
	GetStagingBucketName() string
	GetBucketName() string
	GetClient() *minio.Client
}

type minioManagerImpl struct {
	client            *minio.Client
	stagingBucketName string
	bucketName        string
}

func NewMinioManager(ctx context.Context, config utils.Config) MinioManager {
	// client
	endpoint := fmt.Sprintf("%s:%s", config.MinioHost, config.MinioServerPort)
	accessKeyID := config.MinioRootUser
	secretAccessKey := config.MinioRootPassword

	useSSL := utils.ParseBool(config.MinioUseSSL)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), // TODO: token?
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create MinIO client")
	}

	log.Info().Msg("successfully init MinIO client")

	// buckets
	stagingBucketName := "staging"
	bucketName := config.MinioBucketName
	location := "" // TODO: location?
	opts := minio.MakeBucketOptions{Region: location}

	makeBucket(minioClient, ctx, stagingBucketName, opts)
	makeBucket(minioClient, ctx, bucketName, opts)

	log.Info().Msgf("successfully created bucket: %s, %s", stagingBucketName, bucketName)

	return &minioManagerImpl{
		client:            minioClient,
		stagingBucketName: stagingBucketName,
		bucketName:        bucketName,
	}
}

func (m *minioManagerImpl) GetBucketName() string {
	return m.bucketName
}

func (m *minioManagerImpl) GetStagingBucketName() string {
	return m.stagingBucketName
}

func (m *minioManagerImpl) GetClient() *minio.Client {
	return m.client
}

func makeBucket(client *minio.Client, ctx context.Context, bucketName string, opts minio.MakeBucketOptions) {
	err := client.MakeBucket(ctx, bucketName, opts)
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Info().Err(err).Msgf("bucket %s already exists", bucketName)
		} else {
			log.Fatal().Err(err).Msgf("error when making bucket: %s", bucketName)
		}
	}
}
