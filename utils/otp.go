package utils

import (
	"math/rand"
)

const OtpChars = "1234567890"
const OtpLength = 5

func GenerateOTP() (string, error) {
	buffer := make([]byte, OtpLength)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(OtpChars)
	for i := 0; i < OtpLength; i++ {
		buffer[i] = OtpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
