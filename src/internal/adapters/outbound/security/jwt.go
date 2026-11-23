package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"aman/internal/ports/outbound"
)

var _ outbound.JWTService = (*JWTService)(nil)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) CreateToken(appID, identifier string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: now.Add(ttl).Unix(),
		Issuer:    appID,
		Subject:   identifier,
		IssuedAt:  now.Unix(),
	})
	return claims.SignedString([]byte(os.Getenv("APP_SECRET")))
}

func (s *JWTService) ParseToken(token string) (appID string, identifier string, err error) {
	claims := &jwt.StandardClaims{}
	if _, err = jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("APP_SECRET")), nil
	}); err != nil {
		return "", "", err
	}
	return claims.Issuer, claims.Subject, nil
}
