package usecases

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/daniel-bss/havlabs/internal/media/client"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

type newsUsecaseImpl struct {
	minioManager client.MinioManager
}

func New(m client.MinioManager) NewsUsecase {
	return &newsUsecaseImpl{
		minioManager: m,
	}
}

func (uc *newsUsecaseImpl) CreateUpload(ctx context.Context) (string, error) {
	// TODO: get prefix from grpc req
	objectKey := fmt.Sprintf("images/%s.png", uuid.NewString())

	client := uc.minioManager.GetClient()
	bucketName := uc.minioManager.GetBucketName()
	url1, err := client.PresignedPutObject(ctx, bucketName, objectKey, time.Minute*2)
	if err != nil {
		log.Error().Err(err).Msg("failed when PresignedPutObject")
		return "", err
	}

	o := fmt.Sprintf("images/%s.png", "b9070510-77e9-4d41-b156-f417131818cb")
	x, err := client.StatObject(ctx, bucketName, o, minio.GetObjectOptions{})
	// object_key = article/<user_id>/<media_id> for staging

	/*
		CREATE TYPE media_status AS ENUM (
		  'pending',
		  'uploaded',
		  'ready',
		  'failed'
		);

	*/

	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(">>", x.Metadata)

	reqParams := make(url.Values)
	// reqParams.Add("Content-Type", "image/png")
	url2, err := client.PresignedGetObject(ctx, bucketName, objectKey, time.Minute*2, reqParams)
	if err != nil {
		log.Error().Err(err).Msg("failed when PresignedGetObject")
		return "", err
	}

	return fmt.Sprintf("%s %s", url1.String(), url2), nil
}

/*


You are now in a position to:

implement upload from Next.js

generate public read URLs

store only object keys in DB

add bucket policies

add file size limits

add content-type validation

expire presigned URLs safely

curl -X PUT \
  -H "Content-Type: image/png" \
  --upload-file ./abc.png \
  "http://localhost:9000/mybucket/images/509f8879-cd91-4277-897a-acaf184227bd.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=media%2F20251216%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251216T055224Z&X-Amz-Expires=120&X-Amz-SignedHeaders=host&X-Amz-Signature=6004583ae7074cf2d939bfdadb3d27fe131adfca5d0135d0da044475dba4ef74"


*/
