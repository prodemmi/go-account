package repos

import (
	"time"

	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/session"
	"aman/internal/ports/outbound"
)

var _ outbound.ISessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(s session.Session) (session.Session, error) {
	m := pgmodels.ToModelSession(s)
	if err := r.db.Create(&m).Error; err != nil {
		return session.Session{}, err
	}
	return pgmodels.ToDomainSession(m), nil
}

func (r *SessionRepository) FindAll() ([]session.Session, error) {
	var result []pgmodels.Session
	if err := r.db.Find(&result).Error; err != nil {
		return nil, err
	}
	sessions := make([]session.Session, len(result))
	for i, m := range result {
		sessions[i] = pgmodels.ToDomainSession(m)
	}
	return sessions, nil
}

func (r *SessionRepository) FindByToken(token string) (*session.Session, error) {
	var m pgmodels.Session
	if err := r.db.Where("token = ?", token).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, session.ErrNotFound
		}
		return nil, err
	}
	s := pgmodels.ToDomainSession(m)
	return &s, nil
}

func (r *SessionRepository) FindByAppID(appID string) ([]session.Session, error) {
	var result []pgmodels.Session
	if err := r.db.Where("app_id = ?", appID).Find(&result).Error; err != nil {
		return nil, err
	}
	sessions := make([]session.Session, len(result))
	for i, m := range result {
		sessions[i] = pgmodels.ToDomainSession(m)
	}
	return sessions, nil
}

func (r *SessionRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&pgmodels.Session{}).Error
}

func (r *SessionRepository) DeleteExpired(before time.Time) error {
	return r.db.Where("expires_at < ?", before).Delete(&pgmodels.Session{}).Error
}
