package security

import (
	"errors"
	"strings"
	"time"

	"aman/internal/domains/app"
)

// MockJWTService issues opaque tokens of the form "appID:identifier".
type MockJWTService struct{}

func (m *MockJWTService) CreateToken(appID, identifier string, _ time.Duration) (string, error) {
	return appID + ":" + identifier, nil
}

func (m *MockJWTService) ParseToken(token string) (string, string, error) {
	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid token")
	}
	return parts[0], parts[1], nil
}

// MockAuthHash returns a stable hash for a given app and user agent.
type MockAuthHash struct{}

func (m *MockAuthHash) CreateAuthHash(a app.App, agent string) string {
	return a.ID.String() + "|" + agent
}

func (m *MockAuthHash) VerifyAuthHash(hash string, a app.App, agent string) bool {
	return hash == m.CreateAuthHash(a, agent)
}

// MockPasswordService is a no-op password hasher for tests.
type MockPasswordService struct{}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	return "hashed:" + password, nil
}

func (m *MockPasswordService) CheckPasswordHash(password, hash string) bool {
	return hash == "hashed:"+password
}
