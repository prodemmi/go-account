package otpgen

import "errors"

// MockGenerator returns a fixed code. It is intended for tests.
type MockGenerator struct {
	Code string
	Fail bool
}

func (m *MockGenerator) GenerateOTP(length int, charset string) (string, error) {
	if m.Fail {
		return "", errors.New("otp generation failed")
	}
	if m.Code == "" {
		return "12345", nil
	}
	return m.Code, nil
}
