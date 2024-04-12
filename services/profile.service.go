package services

import (
	"account/models"
	"account/repositories"
)

type (
	ProfileService struct {
		ProfileRepository repositories.IProfileRepository
	}
	IProfileService interface {
		Create(profile models.Profile) (*models.Profile, error)
		GetProfile(id string) (*models.Profile, error)
		GetProfileByUserID(id string) (*models.Profile, error)
	}
)

func (p ProfileService) Create(profile models.Profile) (*models.Profile, error) {
	result, err := p.Create(profile)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p ProfileService) GetProfile(id string) (*models.Profile, error) {
	result, err := p.GetProfile(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p ProfileService) GetProfileByUserID(id string) (*models.Profile, error) {
	result, err := p.GetProfileByUserID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewProfileService(profileRepository repositories.IProfileRepository) IProfileService {
	return ProfileService{
		ProfileRepository: profileRepository,
	}
}
