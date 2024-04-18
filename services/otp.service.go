package services

import (
	"account/models"
	"account/repositories"
	"github.com/google/uuid"
)

type (
	OTPService struct {
		OTPRepository repositories.IOTPRepository
	}
	IOTPService interface {
		CreateOTP(userID uuid.UUID, appID uuid.UUID, otpType string, otpHash string) (models.Otp, error)
		CheckOTP(userID uuid.UUID, appID uuid.UUID, otpCode string, otpType string, otpHash string) (bool, error)
	}
)

func (o OTPService) CreateOTP(userID uuid.UUID, appID uuid.UUID, otpType string, otpHash string) (models.Otp, error) {
	result, errActive := o.OTPRepository.OTPActive(otpHash)
	if errActive == nil {
		return result, nil
	}

	otp, err := o.OTPRepository.Create(userID, appID, otpType, otpHash)
	if err != nil {
		return models.Otp{}, err
	}

	return otp, nil
}

func (o OTPService) CheckOTP(userID uuid.UUID, appID uuid.UUID, otpCode string, otpType string, otpHash string) (bool, error) {
	result, err := o.OTPRepository.CheckOTP(userID, appID, otpCode, otpType, otpHash)
	if err != nil {
		return false, err
	}

	return result, nil
}

func NewOTPService(otpRepository repositories.IOTPRepository) IOTPService {
	return OTPService{
		OTPRepository: otpRepository,
	}
}
