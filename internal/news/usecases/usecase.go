package usecases

import (
	"context"

	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
)

type NewsUsecase interface {
	GetNews(context.Context) []dtos.NewsDto
	CreateNews(context.Context, db.CreateNewsParams) []dtos.NewsDto
}
