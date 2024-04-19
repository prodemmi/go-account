//go:build unit_tests
// +build unit_tests

package utils

import (
	"account/utils"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateToken(t *testing.T) {
	got, err := utils.CreateToken("4462ebfc-5f91-4ef0-9cfb-ac6e7687a66e", "cd613e30-d8f1-4adf-91b7-584a2265b1f5")
	assert.Nil(t, err)
	assert.True(t, strings.HasPrefix(got, "ey"))
}
