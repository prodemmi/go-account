package services

import (
	"github.com/gin-gonic/gin"
	"account/models"
	"gorm.io/gorm"
)

type (
	AppService struct {
		App *gin.Engine
		DB  *gorm.DB
	}
	appService interface {
		Create(name string, url string) (*models.App, error)
		FindAll() ([]*models.App, error)
		FindOne(appID string) (*models.App, error)
		Exists(appID string) (bool, error)
	}
)

func (a AppService) Create(name string, url string) (*models.App, error) {
	var result models.App

	err := a.DB.Model(&result).Create(&models.App{
		Name: name,
		Url:  url,
	}).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a AppService) FindAll() ([]*models.App, error) {
	var result []*models.App

	err := a.DB.Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a AppService) FindOne(appID string) (*models.App, error) {
	var result models.App

	err := a.DB.Where("id = ?", appID).First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a AppService) Exists(appID string) (bool, error) {
	var exists bool
	err := a.DB.Model(&models.App{}).
		Select("count(*) > 0").
		Where("id = ?", appID).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func CreateNewAppService(app *gin.Engine, db *gorm.DB) *AppService {
	return &AppService{
		App: app,
		DB:  db,
	}
}
