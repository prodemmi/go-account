package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const passwordString = "123456"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(passwordString)

	assert.Nil(t, err)
	assert.NotNil(t, hash)
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword(passwordString)
	result := CheckPasswordHash(passwordString, hash)

	assert.True(t, result)
}
