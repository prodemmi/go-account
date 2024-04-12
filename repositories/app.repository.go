package repositories

import (
	"account/internal"
	"account/models"
)

type (
	AppRepository  struct{}
	IAppRepository interface {
		Create(app models.App) (*models.App, error)
		FindAll() ([]models.App, error)
		FindOne(appID string) (*models.App, error)
		Exists(appID string) (bool, error)
		DeleteOne(id string) (bool, error)
	}
)

func (a AppRepository) Create(app models.App) (*models.App, error) {
	err := internal.DB.Create(&app).Error

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (a AppRepository) FindAll() ([]models.App, error) {
	var result []models.App

	err := internal.DB.Find(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a AppRepository) FindOne(appID string) (*models.App, error) {
	var result models.App

	err := internal.DB.Where("id = ?", appID).First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a AppRepository) Exists(appID string) (bool, error) {
	var exists bool
	err := internal.DB.Model(&models.App{}).
		Select("count(*) > 0").
		Where("id = ?", appID).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a AppRepository) DeleteOne(id string) (bool, error) {
	var app models.App
	err := internal.DB.
		Where("id = ?", id).
		Delete(&app).
		Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func NewAppRepository() IAppRepository {
	return AppRepository{}
}
