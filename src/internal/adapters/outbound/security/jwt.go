package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"aman/internal/ports/outbound"
)

var _ outbound.IJWTService = (*JWTService)(nil)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) CreateToken(appID string, sessionTimeout time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(sessionTimeout).Unix(),
		Issuer:    appID,
		Audience:  "guest",
		IssuedAt:  time.Now().Unix(),
	})
	return claims.SignedString([]byte(os.Getenv("APP_SECRET")))
}

func (s *JWTService) ParseToken(token string) (appID string, err error) {
	claims := &jwt.StandardClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("APP_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	if !parsed.Valid {
		return "", errors.New("invalid token")
	}
	if err := claims.Valid(); err != nil {
		return "", err
	}
	if claims.Issuer == "" {
		return "", errors.New("missing issuer")
	}
	return claims.Issuer, nil
}
