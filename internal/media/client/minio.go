package client

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

func NewMinio() *minio.Client {
	endpoint := "localhost:9000"
	accessKeyID := "media"
	secretAccessKey := "mediapwd"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create MinIO client")
	}

	return minioClient
}
