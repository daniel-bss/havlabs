package usecases

type NewsUsecaseImpl struct {
}

func New() NewsUsecase {
	return &NewsUsecaseImpl{}
}
