package integration

import (
	"account/database/models"
)

type testData struct {
	App           models.App
	User          models.User
	OtpHash       string
	OtpCode       string
	AuthUserID    string
	AuthUserToken string
}

func (t *testData) Init() {}

func (t *testData) SetUser(user models.User) {
	t.User = user
}

func (t *testData) GetApp() models.App {
	return t.App
}

func (t *testData) SetApp(app models.App) {
	t.App = app
}

func (t *testData) SetOtpHash(otpHash string) {
	t.OtpHash = otpHash
}

func (t *testData) SetOtpCode(otpCode string) {
	t.OtpCode = otpCode
}

func (t *testData) SetAuthUserToken(token string) {
	t.AuthUserToken = token
}

func (t *testData) Reset() {
	t.AuthUserToken = ""
	t.OtpHash = ""
	t.OtpCode = ""
	t.AuthUserID = ""
	t.App = models.App{}
	t.User = models.User{}
}

var TestData = testData{}
