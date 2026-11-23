package repos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/otp"
	"aman/internal/ports/outbound"
)

var _ outbound.OTPRepository = (*OTPRepository)(nil)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(o otp.Otp) (otp.Otp, error) {
	m := pgmodels.ToModelOtp(o)
	if err := r.db.Create(&m).Error; err != nil {
		return otp.Otp{}, err
	}
	return pgmodels.ToDomainOtp(m), nil
}

func (r *OTPRepository) Consume(appID, identifier, code, hash string) (otp.Otp, error) {
	uid, err := uuid.Parse(appID)
	if err != nil {
		return otp.Otp{}, otp.ErrInvalid
	}

	var m pgmodels.Otp
	result := r.db.
		Where("app_id = ?", uid).
		Where("identifier = ?", identifier).
		Where("code = ?", code).
		Where("used = ?", false).
		Where("hash = ?", hash).
		Where("expire_at > ?", time.Now()).
		First(&m)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return otp.Otp{}, otp.ErrNotFound
		}
		return otp.Otp{}, result.Error
	}

	// Flip the flag conditionally so two concurrent verifications cannot both
	// consume the same code and create two sessions.
	result = r.db.Model(&pgmodels.Otp{}).
		Where("id = ? AND used = ?", m.ID, false).
		Updates(map[string]any{"used": true, "updated_at": time.Now()})
	if result.Error != nil {
		return otp.Otp{}, result.Error
	}
	if result.RowsAffected == 0 {
		return otp.Otp{}, otp.ErrNotFound
	}

	m.Used = true
	return pgmodels.ToDomainOtp(m), nil
}
