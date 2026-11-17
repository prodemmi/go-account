package otpgen

import (
	"crypto/rand"

	"aman/internal/ports/outbound"
)

var _ outbound.IOTPGenerator = (*OTPGenerator)(nil)

type OTPGenerator struct{}

func NewOTPGenerator() *OTPGenerator {
	return &OTPGenerator{}
}

const otpChars = "1234567890"
const otpLen = 5

func (g *OTPGenerator) GenerateOTP() (string, error) {
	buffer := make([]byte, otpLen)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	charsLen := len(otpChars)
	for i := range buffer {
		buffer[i] = otpChars[int(buffer[i])%charsLen]
	}
	return string(buffer), nil
}
