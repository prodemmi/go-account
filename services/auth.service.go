package services

import "account/repositories"

type (
	AuthService struct {
		AuthRepository repositories.IAuthRepository
	}
	IAuthService interface {
	}
)

func NewAuthService(authRepository repositories.IAuthRepository) IAuthService {
	return AuthService{
		AuthRepository: authRepository,
	}
}
