package usecases

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
	"github.com/google/uuid"
)

type NewsUsecaseImpl struct {
	store db.Store
}

func New(store db.Store) NewsUsecase {
	return &NewsUsecaseImpl{
		store,
	}
}

func (uc *NewsUsecaseImpl) GetNews(ctx context.Context) ([]dtos.NewsDto, error) {
	news, err := uc.store.GetAllNews(ctx)
	if err != nil {
		return nil, err
	}

	return dtos.ConvertNewsArray(news), nil
}

func (uc *NewsUsecaseImpl) CreateNews(ctx context.Context, req *pb.CreateNewsRequest) (uuid.UUID, error) {
	arg := db.CreateNewsParams{
		CreatorUsername: req.CreatorUsername,
		Title:           req.Title,
		Content:         req.Content,
	}

	return uc.store.CreateNews(ctx, arg)
}
