package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/dtos"
	"github.com/daniel-bss/havlabs/internal/media/utils"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

func (uc *newsUsecaseImpl) GetMediaById(ctx context.Context, req *pb.GetOneMediaByIdRequest) (*dtos.Media, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid media id")
	}

	media, err := uc.store.GetMediaById(ctx, id)
	if err != nil {
		return nil, err
	}

	imgMetadata := &ImageMetadata{}
	if err = json.Unmarshal(media.Metadata, imgMetadata); err != nil {
		return nil, err
	}

	client := uc.minioManager.GetClient()
	bucket := uc.minioManager.GetBucketName()
	objectKey := fmt.Sprintf("%s/%s.%s", media.Purpose, media.ID, imgMetadata.Format)
	opts := minio.StatObjectOptions{}

	st, err := client.StatObject(ctx, bucket, objectKey, opts)
	if err != nil {
		if err.Error() == "The specified key does not exist." {
			return nil, utils.NewBadRequestError("invalid media id")
		}

		log.Error().Err(err).Msg("error from StatObject")
		return nil, err
	}
	fmt.Println(st)

	return &dtos.Media{
		Id: media.ID,
	}, nil
}
