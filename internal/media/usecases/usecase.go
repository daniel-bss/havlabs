package usecases

import "context"

type NewsUsecase interface {
	CreateUpload(context.Context) (string, error)
}
