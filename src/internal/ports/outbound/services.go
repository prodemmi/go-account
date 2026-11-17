package outbound

import (
	"time"

	"aman/internal/domains/app"
)

// IAuthService is the outbound port for authentication (CAS) operations.
type IAuthService interface {
	// SignIn authenticates app credentials and returns a session token.
	SignIn(appID string, credentials string) (string, error)
	// Check validates a session token and returns the app ID.
	Check(token string) (appID string, err error)
}

// IPasswordService is the outbound port for password hashing and comparison.
type IPasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// IJWTService is the outbound port for JWT issuance and parsing.
type IJWTService interface {
	CreateToken(appID string, sessionTimeout time.Duration) (string, error)
	ParseToken(token string) (appID string, err error)
}

// IAuthHash is the outbound port for the signing/auth hash.
type IAuthHash interface {
	CreateAuthHash(a app.App, agent string) string
	VerifyAuthHash(hash string, a app.App, agent string) bool
}

// IOTPGenerator is the outbound port for generating one-time codes.
type IOTPGenerator interface {
	GenerateOTP() (string, error)
}

// INotificationService is the outbound port for sending notifications (SMS, etc).
type INotificationService interface {
	SendSMS(to, message string) error
}
