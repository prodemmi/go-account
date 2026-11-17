package outbound

import (
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/otp"
	"aman/internal/domains/session"
)

// IAppRepository is the outbound port for app persistence.
type IAppRepository interface {
	Create(a app.App) (app.App, error)
	FindByID(id string) (*app.App, error)
	FindAll() ([]app.App, error)
	Update(a app.App) (app.App, error)
	Delete(id string) error
}

// ISessionRepository is the outbound port for session persistence.
type ISessionRepository interface {
	Create(s session.Session) (session.Session, error)
	FindAll() ([]session.Session, error)
	FindByToken(token string) (*session.Session, error)
	FindByAppID(appID string) ([]session.Session, error)
	Delete(token string) error
	DeleteExpired(before time.Time) error
}

// IOTPRepository is the outbound port for OTP persistence.
type IOTPRepository interface {
	Create(o otp.Otp) (otp.Otp, error)
	CheckOTP(appID string, code, otpType, hash string) (bool, error)
	OTPActive(hash string) (otp.Otp, error)
}
