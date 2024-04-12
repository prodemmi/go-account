package repositories

import (
	"account/internal"
	"account/models"
	"account/utils"
)

type (
	UserRepository  struct{}
	IUserRepository interface {
		Create(user models.User) (*models.User, error)
		FindAll(appID string) ([]models.User, error)
		FindOne(id string) (*models.User, error)
		FindByToken(token string) (*models.User, error)
	}
)

func (u UserRepository) Create(user models.User) (*models.User, error) {
	err := internal.DB.Model(user).Create(&user).Error

	if err != nil {
		return nil, err
	}

	token, errToken := utils.CreateToken(user.AppID.String(), user.ID.String())
	if errToken != nil {
		return nil, errToken
	}

	updateErr := internal.DB.Model(&user).
		Update("api_token", token).
		Error
	if updateErr != nil {
		return nil, updateErr
	}
	return &user, nil
}

func (u UserRepository) FindAll(appID string) ([]models.User, error) {
	var result []models.User

	err := internal.DB.
		Where("app_id = ?", appID).
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u UserRepository) FindOne(id string) (*models.User, error) {
	var result models.User

	err := internal.DB.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u UserRepository) FindByToken(token string) (*models.User, error) {
	var result models.User

	err := internal.DB.
		Where("api_token = ?", token).
		First(&result).
		Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func NewUserRepository() IUserRepository {
	return UserRepository{}
}
