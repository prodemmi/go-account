package utils

import (
	"account/database/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

const hashE = "f099f89fd5c6358bb4f37161a093585f923e3eaa"

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
	hash := CreateAuthHash(getUser())

	assert.IsType(t, "", hash)
	assert.Equal(t, hashE, hash)
}

func TestAuthHashVerify(t *testing.T) {
	verify := VerifyAuthHash(hashE, getUser())

	assert.Equal(t, true, verify)
}
