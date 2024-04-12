package utils

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func CreateToken(appID string, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"e": time.Now().Add(5 * time.Minute).Unix(),
		"a": appID,
		"u": userID,
	})
	tokenString, err := token.SignedString([]byte("config.Auth.JWT.Secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
