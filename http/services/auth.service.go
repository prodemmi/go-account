package services

import (
	"account/database/models"
	"account/http/repositories"
)

type (
	AuthService struct {
		AuthRepository repositories.IAuthRepository
		UserService    IUserService
		AppService     IAppService
	}
	IAuthService interface {
		SignUp(user models.User) (models.User, error)
		SignInWithEmail(email string, password string) (models.User, error)
		SignInWithMobile(email string, password string) (models.User, error)
		Check(app string, user string) (bool, error)
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

func (a AuthService) Check(app string, token string) (bool, error) {
	userExists, err := a.UserService.ExistsByToken(app, token)
	if err != nil {
		return false, err
	}

	if userExists == false {
		return false, nil
	}

	appExists, err := a.AppService.CheckExists(app)
	if err != nil {
		return false, err
	}

	if appExists == false {
		return false, nil
	}

	if userExists == true && appExists == true {
		return true, nil
	}

	return false, nil
}

func NewAuthService(authRepository repositories.IAuthRepository, appService IAppService, userService IUserService) IAuthService {
	return AuthService{
		AuthRepository: authRepository,
		UserService:    userService,
		AppService:     appService,
	}
}
