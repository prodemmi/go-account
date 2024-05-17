package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(appID string, userID string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(96 * time.Hour).Unix(),
		Issuer:    appID,
		Subject:   userID,
		Audience:  "guest",
		IssuedAt:  time.Now().Unix(),
	})
	tokenString, err := claims.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
