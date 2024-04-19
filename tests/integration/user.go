package integration

import (
	"account/database/models"
	"account/http/controllers"
	"account/utils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

type StoreUserRes struct {
	Data   models.User `json:"data"`
	Status string      `json:"status"`
}

type IndexUsersRes struct {
	Data   []models.User `json:"data"`
	Status string        `json:"status"`
}

type DeleteUserRes struct {
	Data   bool   `json:"data"`
	Status string `json:"status"`
}

func UserStoreWithTrueBody(t *testing.T) {
	body := controllers.StoreUserBody{
		FirstName: "Emad",
		LastName:  "Malekpour",
		Mobile:    "+989130243519",
		Email:     "prodemmi@gmail.com",
	}
	req, _ := http.NewRequest("POST", utils.GetTestBaseUrl("users"), utils.StructToJsonBody(body))
	req.Header.Add("APP-ID", TestData.App.ID.String())

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes StoreUserRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		TestData.SetUser(jsonRes.Data)

		assert.Equal(t, body.FirstName, jsonRes.Data.FirstName)
		assert.Equal(t, "success", jsonRes.Status)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}
}

func UserGetUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", utils.GetTestBaseUrl("users"), nil)
	req.Header.Add("APP-ID", TestData.App.ID.String())
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes IndexUsersRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		assert.Equal(t, 1, len(jsonRes.Data))
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}

// TODO:add delete user
//func User_delete_one_user(t *testing.T) {
//	req, _ := http.NewRequest("DELETE", utils.GetTestBaseUrl("users/"+testUser), nil)
//	res, err := http.DefaultClient.Do(req)
//	defer res.Body.Close()
//
//	if err != nil {
//		log.Fatal(err.Error())
//	} else {
//		var jsonRes DeleteUserRes
//
//		json.NewDecoder(res.Body).Decode(&jsonRes)
//
//		assert.True(t, jsonRes.Data)
//		assert.Equal(t, "success", jsonRes.Status)
//		assert.Equal(t, http.StatusOK, res.StatusCode)
//	}
//}
