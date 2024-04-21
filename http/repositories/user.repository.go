package repositories

import (
	"account/database/models"
	"account/internal"
	"fmt"
)

type (
	UserRepository  struct{}
	IUserRepository interface {
		Create(user models.User) (models.User, error)
		FindAll(appID string) ([]models.User, error)
		FindOne(id string) (*models.User, error)
		FindOneBy(where map[string]interface{}) (models.User, error)
		FindByToken(token string) (*models.User, error)
		CheckAuthenticated(email string, mobile string) (bool, error)
		Update(id string, where map[string]interface{}) (models.User, error)
		DeleteToken(appID string, userID string) (bool, error)
	}
)

func (u UserRepository) DeleteToken(appID string, userID string) (bool, error) {
	err := internal.DB.
		Model(&models.User{}).
		Where("app_id = ?", appID).
		Where("id = ?", userID).
		Update("api_token", nil).
		Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (u UserRepository) Create(user models.User) (models.User, error) {
	exists, errExists := u.CheckAuthenticated(*user.Email, *user.Mobile)
	if errExists != nil {
		return models.User{}, errExists
	}

	if exists {
		return models.User{}, fmt.Errorf("user exists please try sign in")
	}
	err := internal.DB.Create(&user).Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u UserRepository) FindAll(appID string) ([]models.User, error) {
	var result []models.User

	err := internal.DB.
		Where("app_id = ?", appID).
		Order("created_at desc").
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

func (u UserRepository) CheckAuthenticated(email string, mobile string) (bool, error) {
	var exists bool
	err := internal.DB.Model(&models.User{}).
		Select("count(*) > 0").
		Where("signed_up", true).
		Where(internal.DB.Where("email = ?", email).Or(internal.DB.Where("mobile = ?", mobile))).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (u UserRepository) Update(id string, where map[string]interface{}) (models.User, error) {
	var result models.User
	err := internal.DB.
		Model(&result).
		Updates(where).
		Error

	if err != nil {
		return models.User{}, err
	}

	return result, nil
}

func (u UserRepository) FindOneBy(where map[string]interface{}) (models.User, error) {
	var user models.User
	err := internal.DB.
		Where(where).
		First(&user).
		Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func NewUserRepository() IUserRepository {
	return UserRepository{}
}
