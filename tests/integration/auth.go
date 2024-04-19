package integration

import (
	"account/http/controllers"
	"account/http/repositories"
	"account/utils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strings"
	"testing"
)

type SignUpDataRes struct {
	UserID  string `json:"user_id"`
	OtpHash string `json:"otp_hash"`
}

type SignUpAuthRes struct {
	Data   SignUpDataRes `json:"data"`
	Status string        `json:"status"`
}

func AuthSignUp(t *testing.T) {
	body := controllers.SignUpBody{
		FirstName: "Ali",
		LastName:  "Akbari",
		Email:     "12345@gmail.com",
		Mobile:    "+989130264515",
		Password:  "123456",
	}

	req, _ := http.NewRequest("POST", utils.GetTestBaseUrl("auth/sign-up"), utils.StructToJsonBody(body))
	req.Header.Add("APP-ID", TestData.App.ID.String())

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes SignUpAuthRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		TestData.AuthUserID = jsonRes.Data.UserID
		TestData.SetOtpHash(jsonRes.Data.OtpHash)

		assert.True(t, strings.Contains(jsonRes.Data.UserID, "-"))
		assert.IsType(t, "", jsonRes.Data.UserID)
		assert.Equal(t, "success", jsonRes.Status)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}

type VerifyAuthDataRes struct {
	ApiToken   string `json:"api_token"`
	RedirectTo string `json:"redirect_to"`
}

type VerifyAuthRes struct {
	Data   VerifyAuthDataRes `json:"data"`
	Status string            `json:"status"`
}

func AuthVerifySignUp(t *testing.T) {
	otp, errOtp := repositories.OTPRepository{}.OTPActive(TestData.OtpHash)
	if errOtp != nil {
		log.Fatal(errOtp.Error())
	}

	body := controllers.VerifyOTPBody{
		OtpHash: TestData.OtpHash,
		OtpType: TestData.App.OtpType,
		UserID:  TestData.AuthUserID,
		OtpCode: otp.Code,
	}

	req, _ := http.NewRequest("POST", utils.GetTestBaseUrl("auth/verify-otp"), utils.StructToJsonBody(body))
	req.Header.Add("APP-ID", TestData.App.ID.String())

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err.Error())
	} else {
		var jsonRes VerifyAuthRes

		_ = json.NewDecoder(res.Body).Decode(&jsonRes)

		TestData.SetAuthUserToken(jsonRes.Data.ApiToken)

		assert.True(t, strings.HasPrefix(jsonRes.Data.ApiToken, "ey"))
		assert.Equal(t, TestData.App.Url, jsonRes.Data.RedirectTo)
		assert.Equal(t, "success", jsonRes.Status)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}
}
