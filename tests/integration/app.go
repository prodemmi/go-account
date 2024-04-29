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

type StoreAppRes struct {
	Data   models.App `json:"data"`
	Status string     `json:"status"`
}

func TestAppStoreAppWithTrueBody(t *testing.T) {
	body := controllers.CreateAppBody{
		Name:    "Test App",
		Url:     "https://test-app.com/",
		AuthUrl: "https://test-app.com/auth/callback",
		OtpType: "email",
	}
	res, err := http.Post(utils.GetTestBaseUrl("apps"), "application/json", utils.StructToJsonBody(body))
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes StoreAppRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		TestData.SetApp(jsonRes.Data)

		assert.Equal(t, "Test App", jsonRes.Data.Name)
		assert.Equal(t, "https://test-app.com/", jsonRes.Data.Url)
		assert.Equal(t, "https://test-app.com/auth/callback", jsonRes.Data.AuthUrl)
		assert.Equal(t, "email", jsonRes.Data.OtpType)
		assert.Equal(t, "success", jsonRes.Status)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}
}

type IndexAppsRes struct {
	Data   []models.App `json:"data"`
	Status string       `json:"status"`
}

func TestAppGetApps(t *testing.T) {
	res, err := http.Get(utils.GetTestBaseUrl("apps"))
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes IndexAppsRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		assert.Equal(t, 1, len(jsonRes.Data))
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}

type SingleAppRes struct {
	Data   models.App `json:"data"`
	Status string     `json:"status"`
}

func TestAppGetOneApp(t *testing.T) {
	appID := TestData.App.ID.String()
	res, err := http.Get(utils.GetTestBaseUrl("apps/" + appID))
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes SingleAppRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		assert.Equal(t, appID, jsonRes.Data.ID.String())
		assert.Equal(t, "success", jsonRes.Status)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}

type DeleteAppRes struct {
	Data   bool   `json:"data"`
	Status string `json:"status"`
}

func TestAppDeleteOneApp(t *testing.T) {
	appID := TestData.App.ID.String()
	req, _ := http.NewRequest("DELETE", utils.GetTestBaseUrl("apps/"+appID), nil)
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes DeleteAppRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		assert.True(t, jsonRes.Data)
		assert.Equal(t, "success", jsonRes.Status)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}
