package repositories

import (
	"account/internal"
	"account/models"
	"account/utils"
	"fmt"
	uuid "github.com/google/uuid"
	"time"
)

type (
	OTPRepository  struct{}
	IOTPRepository interface {
		Create(userID uuid.UUID, appID uuid.UUID, otpType string, otpHash string) (models.Otp, error)
		CheckOTP(userID uuid.UUID, appID uuid.UUID, otpCode string, otpType string, otpHash string) (bool, error)
		OTPActive(otpHash string) (models.Otp, error)
	}
)

func (O OTPRepository) Create(userID uuid.UUID, appID uuid.UUID, otpType string, otpHash string) (models.Otp, error) {
	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return models.Otp{}, err
	}

	expireAt := time.Now().Add(2 * time.Minute)
	//expireAt := time.Now().Add(10 * time.Second)
	otp := models.Otp{
		AppID:    appID,
		UserID:   userID,
		Type:     otpType,
		Code:     otpCode,
		Hash:     &otpHash,
		ExpireAt: &expireAt,
	}

	errCreate := internal.DB.Create(&otp).Error
	if errCreate != nil {
		return models.Otp{}, nil
	}

	return otp, nil
}

func (O OTPRepository) CheckOTP(userID uuid.UUID, appID uuid.UUID, otpCode string, otpType string, otpHash string) (bool, error) {
	var otp models.Otp
	result := internal.DB.
		Where("app_id = ?", appID).
		Where("user_id = ?", userID).
		Where("type = ?", otpType).
		Where("code = ?", otpCode).
		Where("used = ?", false).
		Where("hash = ?", otpHash).
		Where("expire_at > ?", time.Now()).
		First(&otp)
	if result.Error != nil {
		return false, result.Error
	}

	exists := result.RowsAffected > 0

	if exists {
		if otp.ExpireAt.After(time.Now()) {
			otp.Used = true
			internal.DB.Save(&otp)
			return true, nil
		}
		return false, fmt.Errorf("otp code has expired")
	}

	return false, fmt.Errorf("otp not exists")
}

func (O OTPRepository) OTPActive(otpHash string) (models.Otp, error) {
	var otp models.Otp

	err := internal.DB.
		Where("expire_at > ?", time.Now()).
		Where("hash = ?", otpHash).
		Where("used = ?", false).
		First(&otp).
		Error
	if err != nil {
		return models.Otp{}, fmt.Errorf("no active otp exists")
	}
	fmt.Println("otp ", otp)

	return otp, nil
}

func NewOTPRepository() IOTPRepository {
	return OTPRepository{}
}
