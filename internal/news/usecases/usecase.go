package usecases

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/news/client"
	db "github.com/daniel-bss/havlabs/internal/news/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/news/dtos"
	"github.com/google/uuid"
)

type NewsUsecase interface {
	GetNews(context.Context) ([]dtos.NewsDto, error)
	CreateNews(context.Context, *pb.CreateNewsRequest) (*uuid.UUID, error)
}

type newsUsecaseImpl struct {
	store       db.Store
	mediaClient client.MediaClient
}

func New(store db.Store, mediaClient client.MediaClient) NewsUsecase {
	return &newsUsecaseImpl{
		store,
		mediaClient,
	}
}
