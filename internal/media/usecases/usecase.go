package usecases

import (
	"context"

	"github.com/daniel-bss/havlabs-proto/pb"
	"github.com/daniel-bss/havlabs/internal/media/client"
	db "github.com/daniel-bss/havlabs/internal/media/db/sqlc"
	"github.com/daniel-bss/havlabs/internal/media/dtos"
)

type NewsUsecase interface {
	CreateUploadSession(context.Context, int, *pb.CreateUploadSessionRequest) (*dtos.UploadSession, error)
	ConfirmUpload(context.Context, *pb.ConfirmUploadRequest) (*dtos.ConfirmUpload, error)
	GetMediaById(context.Context, *pb.GetOneMediaByIdRequest) (*dtos.Media, error)
	GetMediaURLById(context.Context, int, *pb.GetOneMediaByIdRequest) (string, error)
}

type newsUsecaseImpl struct {
	minioManager client.MinioManager
	store        db.Store
}

func New(m client.MinioManager, s db.Store) NewsUsecase {
	return &newsUsecaseImpl{
		minioManager: m,
		store:        s,
	}
}
