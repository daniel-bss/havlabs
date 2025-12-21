package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

func (uc *newsUsecaseImpl) GetMediaURLById(ctx context.Context, duration int, req *pb.GetOneMediaByIdRequest) (string, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return "", utils.NewBadRequestError("invalid media id")
	}

	media, err := uc.store.GetMediaById(ctx, id)
	if err != nil {
		return "", err
	}

	imgMetadata := &ImageMetadata{}
	if err = json.Unmarshal(media.Metadata, imgMetadata); err != nil {
		return "", err
	}

	client := uc.minioManager.GetClient()
	bucket := uc.minioManager.GetBucketName()
	objectKey := fmt.Sprintf("%s/%s.%s", media.Purpose, media.ID, imgMetadata.Format)
	opts := minio.StatObjectOptions{}

	// check object exists
	_, err = client.StatObject(ctx, bucket, objectKey, opts)
	if err != nil {
		if err.Error() == "The specified key does not exist." {
			return "", utils.NewBadRequestError("invalid media id")
		}

		log.Error().Err(err).Msg("error from StatObject")
		return "", err
	}

	// create presigned GET URL
	imageConcreteURL, err := client.PresignedGetObject(ctx, bucket, objectKey, time.Minute*time.Duration(duration), make(url.Values))
	if err != nil {
		log.Error().Err(err).Msgf("error from PresignedGetObject. media_id %s", req.Id)
		return "", err
	}

	return imageConcreteURL.String(), nil
}
