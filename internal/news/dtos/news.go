package dtos

import (
	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/google/uuid"
)

type NewsDto struct {
	ID              uuid.UUID
	CreatorUsername string
	Title           string
	Content         string
}

func ConvertNews(n db.News) NewsDto {
	return NewsDto{
		ID:              n.ID,
		CreatorUsername: n.CreatorUsername,
		Title:           n.Title,
		Content:         n.Content,
	}
}

func ConvertNewsArray(news []db.News) []NewsDto {
	newsDto := []NewsDto{}
	for _, n := range news {
		newsDto = append(newsDto, ConvertNews(n))
	}

	return newsDto
}
