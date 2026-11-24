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
	if q.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *q.CreatedFrom)
	}
	if q.CreatedUntil != nil {
		query = query.Where("created_at < ?", *q.CreatedUntil)
	}
	now := time.Now()
	if condition, args := sessionStatusFilter(q.Status, now); condition != "" {
		query = query.Where(condition, args...)
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

func (r *SessionRepository) Stats(appID, status string, now time.Time, since, until *time.Time) (outbound.SessionStats, error) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	activeCondition, activeArgs := sessionStatusFilter(session.StatusActive, now)
	pendingCondition, pendingArgs := sessionStatusFilter(session.StatusPending, now)
	expiredCondition, expiredArgs := sessionStatusFilter("expired", now)
	revokedCondition, revokedArgs := sessionStatusFilter(session.StatusRevoked, now)

	where := "deleted_at IS NULL"
	if appID != "" {
		where += " AND app_id = ?"
	}
	if since != nil {
		where += " AND created_at >= ?"
	}
	if until != nil {
		where += " AND created_at < ?"
	}
	statusCondition, statusArgs := sessionStatusFilter(status, now)
	if statusCondition != "" {
		where += " AND " + statusCondition
	}

	query := `SELECT
		count(*) AS total,
		count(*) FILTER (WHERE ` + activeCondition + `) AS active,
		count(*) FILTER (WHERE ` + pendingCondition + `) AS pending,
		count(*) FILTER (WHERE ` + expiredCondition + `) AS expired,
		count(*) FILTER (WHERE ` + revokedCondition + `) AS revoked,
		count(*) FILTER (WHERE infinite = true) AS infinite,
		count(*) FILTER (WHERE created_at >= ?) AS today,
		count(DISTINCT identifier) AS unique
		FROM sessions WHERE ` + where

	args := append(activeArgs, pendingArgs...)
	args = append(args, expiredArgs...)
	args = append(args, revokedArgs...)
	args = append(args, startOfDay)
	if appID != "" {
		args = append(args, appID)
	}
	if since != nil {
		args = append(args, *since)
	}
	if until != nil {
		args = append(args, *until)
	}
	args = append(args, statusArgs...)

	var stats outbound.SessionStats
	if err := r.db.Raw(query, args...).Scan(&stats).Error; err != nil {
		return outbound.SessionStats{}, err
	}
	return stats, nil
}

func (r *SessionRepository) DailyCounts(appID, status string, since time.Time, until *time.Time) ([]outbound.DailyCount, error) {
	query := r.db.Table("sessions").
		Select("date_trunc('day', created_at) AS day, count(*) AS count").
		Where("deleted_at IS NULL AND created_at >= ?", since)
	if until != nil {
		query = query.Where("created_at < ?", *until)
	}
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}
	if condition, args := sessionStatusFilter(status, time.Now()); condition != "" {
		query = query.Where(condition, args...)
	}

	var rows []outbound.DailyCount
	if err := query.Group("date_trunc('day', created_at)").Order("day ASC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// sessionStatusFilter is shared by the session table and dashboard aggregates
// so status categories use the same lifecycle semantics everywhere.
func sessionStatusFilter(status string, now time.Time) (string, []any) {
	switch status {
	case session.StatusActive:
		return legacyActive + " AND expires_at > ?", []any{now}
	case "expired":
		return expired, []any{now}
	case session.StatusPending:
		return "status = ? AND expires_at > ?", []any{session.StatusPending, now}
	case session.StatusRevoked:
		return legacyRevoked, nil
	default:
		return "", nil
	}
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
