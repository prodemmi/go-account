package repos

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/otp"
	"aman/internal/ports/outbound"
)

var _ outbound.IOTPRepository = (*OTPRepository)(nil)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(o otp.Otp) (otp.Otp, error) {
	m := pgmodels.ToModelOtp(o)
	if err := r.db.Create(&m).Error; err != nil {
		return otp.Otp{}, nil
	}
	return pgmodels.ToDomainOtp(m), nil
}

func (r *OTPRepository) CheckOTP(appID string, code, otpType, hash string) (bool, error) {
	uid, err := uuid.Parse(appID)
	if err != nil {
		return false, err
	}

	now := time.Now()
	result := r.db.Model(&pgmodels.Otp{}).
		Where("app_id = ?", uid).
		Where("type = ?", otpType).
		Where("code = ?", code).
		Where("used = ?", false).
		Where("hash = ?", hash).
		Where("expire_at > ?", now).
		Update("used", true)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, fmt.Errorf("otp not exists")
	}
	return true, nil
}

func (r *OTPRepository) OTPActive(hash string) (otp.Otp, error) {
	var m pgmodels.Otp
	err := r.db.
		Where("expire_at > ?", time.Now()).
		Where("hash = ?", hash).
		Where("used = ?", false).
		Order("created_at DESC").
		First(&m).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return otp.Otp{}, fmt.Errorf("no active otp exists")
		}
		return otp.Otp{}, err
	}
	return pgmodels.ToDomainOtp(m), nil
}
