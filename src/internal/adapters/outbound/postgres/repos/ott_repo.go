package repos

import (
	"time"

	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/ott"
	"aman/internal/ports/outbound"
)

var _ outbound.OTTRepository = (*OTTRepository)(nil)

type OTTRepository struct {
	db *gorm.DB
}

func NewOTTRepository(db *gorm.DB) *OTTRepository {
	return &OTTRepository{db: db}
}

func (r *OTTRepository) Create(o ott.OTT) (ott.OTT, error) {
	m := pgmodels.ToModelOTT(o)
	if err := r.db.Create(&m).Error; err != nil {
		return ott.OTT{}, err
	}
	return pgmodels.ToDomainOTT(m), nil
}

func (r *OTTRepository) FindByToken(token string) (*ott.OTT, error) {
	var m pgmodels.OTT
	if err := r.db.Where("token = ?", token).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ott.ErrNotFound
		}
		return nil, err
	}
	o := pgmodels.ToDomainOTT(m)
	return &o, nil
}

// MarkUsed consumes the ticket exactly once. The conditional update makes a
// concurrent replay fail even when both requests read the ticket as unused.
func (r *OTTRepository) MarkUsed(id string) error {
	result := r.db.Model(&pgmodels.OTT{}).
		Where("id = ? AND used = ?", id, false).
		Updates(map[string]any{"used": true, "updated_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ott.ErrUsed
	}
	return nil
}

func (r *OTTRepository) DeleteExpired(before time.Time) error {
	return r.db.Where("expire_at < ?", before).Delete(&pgmodels.OTT{}).Error
}
