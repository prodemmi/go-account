package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	got, err := GenerateOTP()

	assert.Nil(t, err)
	assert.Len(t, got, OtpLength)
	assert.IsType(t, "", got)
}
