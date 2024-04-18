package services

import (
	"account/database/models"
	"account/http/repositories"
)

type (
	AppService struct {
		AppRepository repositories.IAppRepository
	}
	IAppService interface {
		CreateApp(app models.App) (models.App, error)
		GetAppsList() ([]models.App, error)
		GetApp(id string) (*models.App, error)
		CheckExists(id string) (bool, error)
		DeleteApp(id string) (bool, error)
	}
)

func (a AppService) CreateApp(app models.App) (models.App, error) {
	result, err := a.AppRepository.Create(app)
	if err != nil {
		return models.App{}, err
	}

	return result, nil
}

func (a AppService) GetAppsList() ([]models.App, error) {
	apps, err := a.AppRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (a AppService) GetApp(id string) (*models.App, error) {
	app, err := a.AppRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a AppService) CheckExists(id string) (bool, error) {
	exists, err := a.AppRepository.Exists(id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a AppService) DeleteApp(id string) (bool, error) {
	exists, err := a.AppRepository.DeleteOne(id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func NewAppService(appRepository repositories.IAppRepository) IAppService {
	return AppService{
		AppRepository: appRepository,
	}
}
