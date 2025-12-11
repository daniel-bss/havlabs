package usecases

import "github.com/daniel-bss/havlabs/internal/news/dtos"

type NewsUsecase interface {
	GetNews() []dtos.NewsDto
}

type NewsUsecaseImpl struct {
}

func New() NewsUsecase {
	return &NewsUsecaseImpl{}
}

func (uc *NewsUsecaseImpl) GetNews() []dtos.NewsDto {
	return []dtos.NewsDto{}
}
