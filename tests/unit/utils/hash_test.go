//go:build unit_tests
// +build unit_tests

package utils

import (
	"account/database/models"
	"account/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var hashE = ""

const userAgent = "Mozilla/5.0 (Windows 98; st-ZA; rv:1.9.1.20) Gecko/2011-06-14 07:44:52 Firefox/3.8"

func getUser() models.User {
	appId, _ := uuid.Parse("cd613e30-d8f1-4adf-91b7-584a2265b1f5")
	userID, _ := uuid.Parse("4462ebfc-5f91-4ef0-9cfb-ac6e7687a66e")
	user := models.User{
		FirstName: "Emad",
		AppID:     appId,
	}
	user.ID = userID

	return user
}

func TestAuthHashCreates(t *testing.T) {
	hashE = utils.CreateAuthHash(getUser(), userAgent)

	assert.IsType(t, "", hashE)
}

func TestAuthHashVerify(t *testing.T) {
	verify := utils.VerifyAuthHash(hashE, getUser(), userAgent)

	assert.Equal(t, true, verify)
}
