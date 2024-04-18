package utils

import (
	"account/database/models"
	"crypto/sha1"
	"encoding/hex"
)

func CreateAuthHash(user models.User) string {
	h := sha1.New()
	h.Write([]byte(user.ID.String() + user.FirstName + user.AppID.String()))

	return hex.EncodeToString(h.Sum(nil))
}

func VerifyAuthHash(hash string, user models.User) bool {
	return hash == CreateAuthHash(user)
}
