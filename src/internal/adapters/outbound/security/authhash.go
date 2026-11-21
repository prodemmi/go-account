package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"aman/internal/domains/app"
	"aman/internal/ports/outbound"
)

var _ outbound.IAuthHash = (*AuthHash)(nil)

type AuthHash struct{}

func NewAuthHash() *AuthHash {
	return &AuthHash{}
}

func (s *AuthHash) CreateAuthHash(a app.App, agent string) string {
	mac := hmac.New(sha256.New, []byte(os.Getenv("APP_SECRET")))
	mac.Write([]byte(a.ID.String() + a.Name + agent))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *AuthHash) VerifyAuthHash(hash string, a app.App, agent string) bool {
	return hash == s.CreateAuthHash(a, agent)
}
