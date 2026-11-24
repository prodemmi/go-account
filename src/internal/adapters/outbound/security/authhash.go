package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	"aman/internal/domains/app"
	"aman/internal/ports/outbound"
)

var _ outbound.AuthHash = (*AuthHash)(nil)

type AuthHash struct{}

type otpChallenge struct {
	AppID      string `json:"app_id"`
	Identifier string `json:"identifier"`
	ExpiresAt  int64  `json:"expires_at"`
}

func NewAuthHash() *AuthHash {
	return &AuthHash{}
}

func (s *AuthHash) CreateAuthHash(a app.App, identifier, agent string, expiresAt time.Time) string {
	secret := os.Getenv("APP_SECRET")
	if secret == "" {
		return ""
	}
	payload, err := json.Marshal(otpChallenge{
		AppID:      a.ID.String(),
		Identifier: identifier,
		ExpiresAt:  expiresAt.UnixNano(),
	})
	if err != nil {
		return ""
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return ""
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return ""
	}
	sealed := aead.Seal(nonce, nonce, payload, []byte(a.ID.String()+"\x00"+agent))
	return base64.RawURLEncoding.EncodeToString(sealed)
}

func (s *AuthHash) VerifyAuthHash(hash string, a app.App, identifier, agent string) bool {
	secret := os.Getenv("APP_SECRET")
	if secret == "" || hash == "" || len(hash) > 4096 {
		return false
	}
	sealed, err := base64.RawURLEncoding.DecodeString(hash)
	if err != nil {
		return false
	}
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return false
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return false
	}
	if len(sealed) < aead.NonceSize()+aead.Overhead() {
		return false
	}
	nonce, ciphertext := sealed[:aead.NonceSize()], sealed[aead.NonceSize():]
	payloadBytes, err := aead.Open(nil, nonce, ciphertext, []byte(a.ID.String()+"\x00"+agent))
	if err != nil {
		return false
	}
	var payload otpChallenge
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return false
	}
	return payload.AppID == a.ID.String() && payload.Identifier == identifier && time.Unix(0, payload.ExpiresAt).After(time.Now())
}
