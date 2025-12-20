package usecases

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/utils"
	"github.com/google/uuid"
)

func (uc *newsUsecaseImpl) CreateNews(ctx context.Context, req *pb.CreateNewsRequest) (*uuid.UUID, error) {
	// TODO: handle title already exists

	username, err := utils.ExtractMetadataFromContextWithKey("x-username", ctx)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid credentials")
	}

	mediaId, err := uuid.Parse(req.MediaId)
	if err != nil {
		return nil, utils.NewBadRequestError("invalid media id")
	}

	mediaIdPtr, err := uc.mediaClient.GetMediaById(ctx, mediaId)
	if err != nil {
		return nil, err
	}

	arg := db.CreateNewsParams{
		CreatorUsername: username,
		Title:           req.Title,
		Content:         req.Content,
		MediaID:         *mediaIdPtr,
	}

	id, err := uc.store.CreateNews(ctx, arg)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
