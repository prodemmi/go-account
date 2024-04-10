package services

import (
	"github.com/gin-gonic/gin"
	"account/models"
	"gorm.io/gorm"
)

type (
	UserService struct {
		App          *gin.Engine
		DB           *gorm.DB
		TokenService TokenService
	}
	userService interface {
		Create(user *models.User) (*models.User, error)
		FindAll() ([]models.User, error)
		FindOne(id string) (*models.User, error)
	}
)

func (a UserService) Create(user *models.User) (*models.User, error) {
	err := a.DB.Create(&user).Error

	if err != nil {
		return nil, err
	}

	_, err2 := a.TokenService.Create(user.AppID, user.ID)

	if err2 != nil {
		return nil, err2
	}

	return user, nil
}

func (a UserService) FindAll() ([]models.User, error) {
	var result []models.User

	err := a.DB.
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a UserService) FindOne(id string) (*models.User, error) {
	var result models.User

	err := a.DB.
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}
