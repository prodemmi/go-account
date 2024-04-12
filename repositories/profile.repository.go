package repositories

import (
	"account/internal"
	"account/models"
)

type (
	ProfileRepository  struct{}
	IProfileRepository interface {
		Create(profile models.Profile) (*models.Profile, error)
		FindByID(id string) (*models.Profile, error)
		FindByUserID(id string) (*models.Profile, error)
	}
)

func (p ProfileRepository) Create(profile models.Profile) (*models.Profile, error) {
	err := internal.DB.Create(&profile).Error

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (p ProfileRepository) FindByID(id string) (*models.Profile, error) {
	var result models.Profile
	err := internal.DB.Where("id = ?", id).First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p ProfileRepository) FindByUserID(id string) (*models.Profile, error) {
	var result models.Profile
	err := internal.DB.Where("user_id = ?", id).First(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func NewProfileRepository() IProfileRepository {
	return ProfileRepository{}
}
