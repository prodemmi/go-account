package security

import (
	"crypto/sha1"
	"encoding/hex"

	"aman/internal/domains/app"
	"aman/internal/ports/outbound"
)

var _ outbound.IAuthHash = (*AuthHash)(nil)

type AuthHash struct{}

func NewAuthHash() *AuthHash {
	return &AuthHash{}
}

func (s *AuthHash) CreateAuthHash(a app.App, agent string) string {
	h := sha1.New()
	h.Write([]byte(a.ID.String() + a.Name + agent))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *AuthHash) VerifyAuthHash(hash string, a app.App, agent string) bool {
	return hash == s.CreateAuthHash(a, agent)
}
