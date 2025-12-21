package usecases

import (
	"context"

	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"github.com/rs/zerolog/log"
)

func (uc *newsUsecaseImpl) GetNews(ctx context.Context) ([]dtos.NewsDto, error) {
	news, err := uc.store.GetAllNews(ctx)
	if err != nil {
		return nil, err
	}

	okFromProvider := true
	result := utils.Map(news, func(n db.News) dtos.NewsDto {
		// basic properties
		result := dtos.ConvertNews(n)

		// call to media service for imageUrl
		mediaId := n.MediaID
		imageUrl, err := uc.mediaClient.GetMediaUrlString(ctx, mediaId)
		result.ImageURL = imageUrl

		if err != nil {
			log.Error().Err(err).Msgf("error from minIO to get media URL. media_id %s", mediaId)
			okFromProvider = false
		}

		return result
	})
	if !okFromProvider {
		return nil, utils.ProviderError("error from minio")
	}

	return result, nil
}
