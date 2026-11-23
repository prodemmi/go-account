package outbound

import (
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/otp"
	"aman/internal/domains/ott"
	"aman/internal/domains/session"
)

// AppRepository is the outbound port for app persistence.
type AppRepository interface {
	Create(a app.App) (app.App, error)
	FindByID(id string) (*app.App, error)
	FindAll() ([]app.App, error)
	FindPage(search string, limit, offset int) ([]app.App, int64, error)
	Update(a app.App) (app.App, error)
	Delete(id string) error
}

// SessionQuery is the persistence-level filter for admin session tables.
type SessionQuery struct {
	AppID  string
	Search string
	Status string
	Limit  int
	Offset int
}

// SessionStats is the aggregate used by the dashboard.
type SessionStats struct {
	Total    int64
	Active   int64
	Pending  int64
	Expired  int64
	Revoked  int64
	Infinite int64
	Today    int64
	Unique   int64
}

// DailyCount is one point on the dashboard chart.
type DailyCount struct {
	Day   time.Time
	Count int64
}

// SessionRepository is the outbound port for session persistence.
type SessionRepository interface {
	Create(s session.Session) (session.Session, error)
	FindByID(id string) (*session.Session, error)
	FindByToken(token string) (*session.Session, error)
	FindAll() ([]session.Session, error)
	FindPage(q SessionQuery) ([]session.Session, int64, error)
	FindActiveByAppAndIdentifier(appID, identifier string) ([]session.Session, error)
	LatestCreatedAt(appID, identifier string) (*time.Time, error)
	Stats(appID string, now time.Time, since *time.Time) (SessionStats, error)
	DailyCounts(appID string, since time.Time) ([]DailyCount, error)
	MarkActive(id string) error
	Revoke(id string) error
	RevokeMany(ids []string) error
	RevokeAll(appID, identifier string) error
	DeleteExpired(before time.Time) error
}

// OTPRepository is the outbound port for OTP persistence.
type OTPRepository interface {
	Create(o otp.Otp) (otp.Otp, error)
	// Consume finds a matching, unused, unexpired OTP and marks it used.
	Consume(appID, identifier, code, hash string) (otp.Otp, error)
}

// OTTRepository is the outbound port for one-time ticket persistence.
type OTTRepository interface {
	Create(o ott.OTT) (ott.OTT, error)
	FindByToken(token string) (*ott.OTT, error)
	MarkUsed(id string) error
	DeleteExpired(before time.Time) error
}
