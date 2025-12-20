package usecases

import (
	"context"

	"github.com/daniel-bss/havlabs/internal/news/dtos"
)

func (uc *newsUsecaseImpl) GetNews(ctx context.Context) ([]dtos.NewsDto, error) {
	news, err := uc.store.GetAllNews(ctx)
	if err != nil {
		return nil, err
	}

	return dtos.ConvertNewsArray(news), nil
}
