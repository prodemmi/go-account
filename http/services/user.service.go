package services

import (
	"account/database/models"
	"account/http/repositories"
	"account/internal"
	"account/utils"
)

type (
	UserService struct {
		UserRepository repositories.IUserRepository
	}
	IUserService interface {
		Create(user models.User) (models.User, error)
		GetUsersList(appID string) ([]models.User, error)
		CreateUserToken(id string) (string, error)
		GetUser(id string) (*models.User, error)
		GetUserWithEmail(email string) (models.User, error)
		GetUserWithMobile(mobile string) (models.User, error)
		GetUserByToken(token string) (*models.User, error)
	}
)

func (a UserService) Create(user models.User) (models.User, error) {
	createdUser, err := a.UserRepository.Create(user)
	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (a UserService) CreateUserToken(id string) (string, error) {
	user, err := a.GetUser(id)
	if err != nil {
		return "", err
	}

	token, errToken := utils.CreateToken(user.AppID.String(), user.ID.String())
	if errToken != nil {
		return "", errToken
	}

	user.ApiToken = &token
	user.SignedUp = true
	errSave := internal.DB.Save(&user).Error
	if errSave != nil {
		return "", errSave
	}

	return token, nil
}

func (a UserService) GetUsersList(appID string) ([]models.User, error) {
	users, err := a.UserRepository.FindAll(appID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a UserService) GetUser(id string) (*models.User, error) {
	users, err := a.UserRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (a UserService) GetUserByToken(token string) (*models.User, error) {
	user, err := a.UserRepository.FindByToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a UserService) GetUserWithEmail(email string) (models.User, error) {
	user, err := a.UserRepository.FindOneBy(map[string]interface{}{
		"email": email,
	})

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (a UserService) GetUserWithMobile(mobile string) (models.User, error) {
	user, err := a.UserRepository.FindOneBy(map[string]interface{}{
		"mobile": mobile,
	})

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func NewUserService(userRepository repositories.IUserRepository) IUserService {
	return UserService{
		UserRepository: userRepository,
	}
}
