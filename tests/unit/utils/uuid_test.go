//go:build unit_tests
// +build unit_tests

package utils

import (
	"account/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsUUID(t *testing.T) {
	result := utils.IsUUID("4462ebfc-5f91-4ef0-9cfb-ac6e7687a66e")
	assert.True(t, result)
}
