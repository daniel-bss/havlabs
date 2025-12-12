package usecases

import (
	"context"
	"fmt"

	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
)

type NewsUsecaseImpl struct {
	store db.Store
}

func New(store db.Store) NewsUsecase {
	return &NewsUsecaseImpl{
		store,
	}
}

func (uc *NewsUsecaseImpl) GetNews(ctx context.Context) []dtos.NewsDto {
	return []dtos.NewsDto{}
}

func (uc *NewsUsecaseImpl) CreateNews(ctx context.Context, arg db.CreateNewsParams) []dtos.NewsDto {
	id, err := uc.store.CreateNews(ctx, arg)
	if err != nil {
		// return id, err
		fmt.Println("wah", err)
		return []dtos.NewsDto{}
	}
	fmt.Println("mantap", id)
	return []dtos.NewsDto{}
}
