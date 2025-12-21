package usecases

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func (uc *newsUsecaseImpl) GetNews(ctx context.Context, req *pb.ListNewsRequest) ([]dtos.NewsDto, uint32, error) {

	limit := req.Limit
	page := req.Page

	including := req.Including

	startTime := req.StartTime
	endTime := req.EndTime

	sortBy := req.SortBy.String()
	ascending := req.Ascending

	if limit < 1 || page < 1 {
		return nil, 0, utils.NewBadRequestError("invalid pagination params")
	}

	v, ok := utils.ValidateSortBy(sortBy)
	if !ok {
		return nil, 0, utils.NewBadRequestError("invalid sort_by params")
	}

	// list news
	news, tot, err := uc.store.ListNews(ctx, db.ListNewsParams{
		Limit:     limit,
		Offset:    (page - 1) * limit,
		Including: including,
		StartTime: pgtype.Timestamp{Time: startTime.AsTime(), Valid: startTime.IsValid()},
		EndTime:   pgtype.Timestamp{Time: endTime.AsTime(), Valid: endTime.IsValid()},
		SortBy:    db.SortByTypes(v),
		Ascending: ascending,
	})
	if err != nil {
		return nil, 0, err
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
		return nil, 0, utils.ProviderError("error from minio")
	}

	return result, tot, nil
}
