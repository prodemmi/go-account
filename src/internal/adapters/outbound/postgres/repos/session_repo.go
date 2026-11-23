package repos

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	pgmodels "aman/internal/adapters/outbound/postgres/models"
	"aman/internal/domains/session"
	"aman/internal/ports/outbound"
)

var _ outbound.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// legacyActive matches rows written before the status column existed.
const legacyActive = "(status = 'active' OR (status = '' AND is_logged_in = true))"

// legacyRevoked matches the revoked rows of both storage generations.
const legacyRevoked = "(status = 'revoked' OR (status = '' AND is_logged_in = false))"

// expired matches every session that has not been revoked and whose expiry has
// passed. Pending sessions that were never exchanged for a ticket expire too,
// so they must show up as expired instead of staying "pending" forever.
const expired = "NOT " + legacyRevoked + " AND expires_at <= ?"

func (r *SessionRepository) Create(s session.Session) (session.Session, error) {
	m := pgmodels.ToModelSession(s)
	if m.Status == "" {
		m.Status = session.StatusPending
	}
	if err := r.db.Create(&m).Error; err != nil {
		return session.Session{}, err
	}
	return pgmodels.ToDomainSession(m), nil
}

func (r *SessionRepository) FindByID(id string) (*session.Session, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, session.ErrNotFound
	}

	var m pgmodels.Session
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, session.ErrNotFound
		}
		return nil, err
	}
	s := pgmodels.ToDomainSession(m)
	return &s, nil
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

func (r *SessionRepository) FindAll() ([]session.Session, error) {
	var result []pgmodels.Session
	if err := r.db.Order("created_at DESC").Find(&result).Error; err != nil {
		return nil, err
	}
	sessions := make([]session.Session, len(result))
	for i, m := range result {
		sessions[i] = pgmodels.ToDomainSession(m)
	}
	return sessions, nil
}

func (r *SessionRepository) FindPage(q outbound.SessionQuery) ([]session.Session, int64, error) {
	query := r.db.Model(&pgmodels.Session{})
	if q.AppID != "" {
		query = query.Where("app_id = ?", q.AppID)
	}
	if q.Search != "" {
		like := "%" + q.Search + "%"
		query = query.Where("(identifier ILIKE ? OR ip ILIKE ? OR user_agent ILIKE ? OR city ILIKE ?)", like, like, like, like)
	}
	now := time.Now()
	switch q.Status {
	case session.StatusActive:
		query = query.Where(legacyActive).Where("expires_at > ?", now)
	case "expired":
		query = query.Where(expired, now)
	case session.StatusPending:
		query = query.Where("status = ?", session.StatusPending).Where("expires_at > ?", now)
	case session.StatusRevoked:
		query = query.Where(legacyRevoked)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var result []pgmodels.Session
	if err := query.Order("created_at DESC").Limit(q.Limit).Offset(q.Offset).Find(&result).Error; err != nil {
		return nil, 0, err
	}
	sessions := make([]session.Session, len(result))
	for i, m := range result {
		sessions[i] = pgmodels.ToDomainSession(m)
	}
	return sessions, total, nil
}

func (r *SessionRepository) FindActiveByAppAndIdentifier(appID, identifier string) ([]session.Session, error) {
	var result []pgmodels.Session
	if err := r.db.
		Where("app_id = ?", appID).
		Where("identifier = ?", identifier).
		Where(legacyActive).
		Where("expires_at > ?", time.Now()).
		Order("created_at DESC").
		Find(&result).Error; err != nil {
		return nil, err
	}
	sessions := make([]session.Session, len(result))
	for i, m := range result {
		sessions[i] = pgmodels.ToDomainSession(m)
	}
	return sessions, nil
}

func (r *SessionRepository) LatestCreatedAt(appID, identifier string) (*time.Time, error) {
	var created *time.Time
	err := r.db.Model(&pgmodels.Session{}).
		Select("created_at").
		Where("app_id = ? AND identifier = ?", appID, identifier).
		Order("created_at DESC").
		Limit(1).
		Scan(&created).Error
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *SessionRepository) Stats(appID string, now time.Time, since *time.Time) (outbound.SessionStats, error) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	where := "deleted_at IS NULL"
	if appID != "" {
		where += " AND app_id = ?"
	}
	if since != nil {
		where += " AND created_at >= ?"
	}

	query := `SELECT
		count(*) AS total,
		count(*) FILTER (WHERE ` + legacyActive + ` AND expires_at > ?) AS active,
		count(*) FILTER (WHERE status = 'pending' AND expires_at > ?) AS pending,
		count(*) FILTER (WHERE ` + expired + `) AS expired,
		count(*) FILTER (WHERE ` + legacyRevoked + `) AS revoked,
		count(*) FILTER (WHERE infinite = true) AS infinite,
		count(*) FILTER (WHERE created_at >= ?) AS today,
		count(DISTINCT identifier) AS unique
		FROM sessions WHERE ` + where

	args := []any{now, now, now, startOfDay}
	if appID != "" {
		args = append(args, appID)
	}
	if since != nil {
		args = append(args, *since)
	}

	var stats outbound.SessionStats
	if err := r.db.Raw(query, args...).Scan(&stats).Error; err != nil {
		return outbound.SessionStats{}, err
	}
	return stats, nil
}

func (r *SessionRepository) DailyCounts(appID string, since time.Time) ([]outbound.DailyCount, error) {
	query := r.db.Table("sessions").
		Select("date_trunc('day', created_at) AS day, count(*) AS count").
		Where("deleted_at IS NULL AND created_at >= ?", since)
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}

	var rows []outbound.DailyCount
	if err := query.Group("date_trunc('day', created_at)").Order("day ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SessionRepository) MarkActive(id string) error {
	return r.db.Model(&pgmodels.Session{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       session.StatusActive,
			"is_logged_in": true,
			"updated_at":   time.Now(),
		}).Error
}

func (r *SessionRepository) Revoke(id string) error {
	return r.db.Model(&pgmodels.Session{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       session.StatusRevoked,
			"is_logged_in": false,
			"updated_at":   time.Now(),
		}).Error
}

func (r *SessionRepository) RevokeMany(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Model(&pgmodels.Session{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"status":       session.StatusRevoked,
			"is_logged_in": false,
			"updated_at":   time.Now(),
		}).Error
}

func (r *SessionRepository) RevokeAll(appID, identifier string) error {
	return r.db.Model(&pgmodels.Session{}).
		Where("app_id = ? AND identifier = ?", appID, identifier).
		Updates(map[string]any{
			"status":       session.StatusRevoked,
			"is_logged_in": false,
			"updated_at":   time.Now(),
		}).Error
}

func (r *SessionRepository) DeleteExpired(before time.Time) error {
	return r.db.Where("expires_at < ?", before).Delete(&pgmodels.Session{}).Error
}
