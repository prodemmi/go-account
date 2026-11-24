package outbound

import (
	"time"

	"aman/internal/domains/app"
)

// JWTService is the outbound port implemented by the security adapter.
type JWTService interface {
	CreateToken(appID, identifier string, ttl time.Duration) (string, error)
	ParseToken(token string) (appID string, identifier string, err error)
}

// AuthHash creates and verifies authenticated OTP challenge tokens.
type AuthHash interface {
	CreateAuthHash(a app.App, identifier, agent string, expiresAt time.Time) string
	VerifyAuthHash(hash string, a app.App, identifier, agent string) bool
}

// PasswordService is the outbound port for password hashing and comparison.
type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}
