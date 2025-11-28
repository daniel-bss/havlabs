package usecases

type AuthServiceImpl struct {
}

type AuthService interface {
}

func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}
