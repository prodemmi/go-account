package utils

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func CreateToken(appID string, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		Id:        userID,
		Subject:   appID,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
