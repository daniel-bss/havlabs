package dtos

import (
	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type NewsDto struct {
	ID          uuid.UUID
	Title       string
	Content     string
	ImageURL    string
	PublishedAt pgtype.Timestamp
}

func ConvertNews(n db.News) NewsDto {
	return NewsDto{
		ID:          n.ID,
		Title:       n.Title,
		Content:     n.Content,
		ImageURL:    "test",
		PublishedAt: n.CreatedAt,
	}
}

func ConvertNewsArray(news []db.News) []NewsDto {
	newsDto := []NewsDto{}
	for _, n := range news {
		newsDto = append(newsDto, ConvertNews(n))
	}

	return newsDto
}
