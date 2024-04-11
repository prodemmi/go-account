package services

import (
	"account/models"
	"account/repositories"
)

type (
	UserService struct {
		UserRepository repositories.IUserRepository
	}
	IUserService interface {
		Create(user *models.User) (*models.User, error)
		GetUsersList(appID string) ([]models.User, error)
		GetUser(id string) (*models.User, error)
		GetUserByToken(token string) (*models.User, error)
	}
)

func (a UserService) Create(user *models.User) (*models.User, error) {
	user, err := a.UserRepository.Create(*user)
	if err != nil {
		return nil, err
	}

	return user, nil
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

func NewUserService(userRepository repositories.IUserRepository) IUserService {
	return UserService{
		UserRepository: userRepository,
	}
}
