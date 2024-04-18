package services

import (
	"account/database/models"
	"account/http/repositories"
)

type (
	AuthService struct {
		AuthRepository repositories.IAuthRepository
		UserService    IUserService
	}
	IAuthService interface {
		SignUp(user models.User) (models.User, error)
		SignInWithEmail(email string, password string) (models.User, error)
		SignInWithMobile(email string, password string) (models.User, error)
	}
)

func (a AuthService) SignUp(user models.User) (models.User, error) {
	authUser, errCreate := a.UserService.Create(user)
	if errCreate != nil {
		return models.User{}, errCreate
	}

	return authUser, nil
}

func (a AuthService) SignInWithEmail(email string, password string) (models.User, error) {
	authUser, errCreate := a.UserService.GetUserWithEmail(email)
	if errCreate != nil {
		return models.User{}, errCreate
	}

	return authUser, nil
}

func (a AuthService) SignInWithMobile(mobile string, password string) (models.User, error) {
	authUser, errCreate := a.UserService.GetUserWithMobile(mobile)
	if errCreate != nil {
		return models.User{}, errCreate
	}

	return authUser, nil
}

func NewAuthService(authRepository repositories.IAuthRepository, userService IUserService) IAuthService {
	return AuthService{
		AuthRepository: authRepository,
		UserService:    userService,
	}
}
