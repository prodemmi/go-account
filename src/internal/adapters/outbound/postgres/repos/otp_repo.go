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

	var m pgmodels.Otp
	result := r.db.
		Where("app_id = ?", uid).
		Where("type = ?", otpType).
		Where("code = ?", code).
		Where("used = ?", false).
		Where("hash = ?", hash).
		Where("expire_at > ?", time.Now()).
		First(&m)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, fmt.Errorf("otp not exists")
	}

	if m.ExpireAt.After(time.Now()) {
		m.Used = true
		if err := r.db.Save(&m).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	return false, fmt.Errorf("otp code has expired")
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
