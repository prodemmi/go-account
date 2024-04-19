//go:build unit_tests
// +build unit_tests

package utils

import (
	"account/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

const passwordString = "123456"

func TestHashPassword(t *testing.T) {
	hash, err := utils.HashPassword(passwordString)

	assert.Nil(t, err)
	assert.NotNil(t, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := utils.HashPassword(passwordString)
	result := utils.CheckPasswordHash(passwordString, hash)

	assert.True(t, result)
}
