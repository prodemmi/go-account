package models

// App model info
// @Description App information
type App struct {
	BaseModel

	// name
	Name string `json:"name" example:"Fotomi"`
	// url
	Url string `json:"url" example:"https://fotomi.net/"`
	// url
	OtpType string `json:"otp_type" gorm:"email" example:"email"`
}
