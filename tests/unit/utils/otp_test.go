//go:build unit_tests
// +build unit_tests

package utils

import (
	"account/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	got, err := utils.GenerateOTP()

	assert.Nil(t, err)
	assert.Len(t, got, utils.OtpLength)
	assert.IsType(t, "", got)
}
