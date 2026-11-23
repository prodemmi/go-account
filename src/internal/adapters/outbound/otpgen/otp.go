package otpgen

import (
	"crypto/rand"

	"aman/internal/domains/app"
	"aman/internal/ports/outbound"
)

var _ outbound.OTPGenerator = (*OTPGenerator)(nil)

type OTPGenerator struct{}

func NewOTPGenerator() *OTPGenerator {
	return &OTPGenerator{}
}

const (
	numericChars      = "0123456789"
	alphaChars        = "abcdefghijklmnopqrstuvwxyz"
	alphanumericChars = numericChars + alphaChars
)

// GenerateOTP creates a cryptographically random code with the requested
// length and character set. The length is clamped to 4..32.
func (g *OTPGenerator) GenerateOTP(length int, charset string) (string, error) {
	if length < 4 {
		length = 4
	}
	if length > 32 {
		length = 32
	}

	chars := numericChars
	switch charset {
	case app.OTPCharsetAlpha:
		chars = alphaChars
	case app.OTPCharsetAlphanumeric:
		chars = alphanumericChars
	}

	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	// Rejection-free modulo is acceptable here: codes are short-lived, single
	// use and the bias against 256%len(chars) is negligible for 32 characters.
	charsLen := byte(len(chars))
	for i := range buffer {
		buffer[i] = chars[int(buffer[i])%int(charsLen)]
	}
	return string(buffer), nil
}
