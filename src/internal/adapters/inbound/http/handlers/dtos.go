package handlers

// Request / response DTOs for the HTTP inbound adapter (CAS).

type signUpBody struct {
	FirstName string `json:"first_name" binding:"required" example:"Emad"`
	LastName  string `json:"last_name" binding:"required" example:"Malekpour"`
	Email     string `json:"email" example:"prodemmi@gmail.com"`
	Mobile    string `json:"mobile" example:"+989130243519"`
	Password  string `json:"password" binding:"required" example:"secure"`
}

type signInBody struct {
	Email    string `json:"email" example:"prodemmi@gmail.com"`
	Mobile   string `json:"mobile" example:"+989130243519"`
	Password string `json:"password" binding:"required" example:"secure"`
}

type verifyOTPBody struct {
	UserID  string `json:"user_id" binding:"required" example:"21636369-8b52-4b4a-97b7-50923ceb3ffd"`
	OtpCode string `json:"otp_code" binding:"required" example:"12345"`
	OtpType string `json:"otp_type" binding:"required" example:"email"`
	OtpHash string `json:"otp_hash" binding:"required" example:"fabecfb8ad5895a6a6309e5a7d1b5c48601d1f9e"`
}

type createAppBody struct {
	Name           string `json:"name" binding:"required" example:"fotomi"`
	Url            string `json:"url" binding:"required" example:"https://fotomi.net/"`
	AuthUrl        string `json:"auth_url" binding:"required" example:"https://fotomi.net/auth/callback"`
	OtpType        string `json:"otp_type" binding:"required" example:"email"`
	SessionTimeout int    `json:"session_timeout" binding:"required" example:"3600"`
}

type updateAppBody struct {
	Name           string `json:"name" example:"fotomi"`
	Url            string `json:"url" example:"https://fotomi.net/"`
	AuthUrl        string `json:"auth_url" example:"https://fotomi.net/auth/callback"`
	OtpType        string `json:"otp_type" example:"email"`
	SessionTimeout int    `json:"session_timeout" example:"3600"`
}

type createSessionBody struct {
	AppID string `json:"app_id" binding:"required" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	Token string `json:"token" binding:"required" example:"ticket-ST-..."`
}

type sendSMS struct {
	To      string `json:"to" binding:"required" example:"+989130243519"`
	Message string `json:"message" binding:"required" example:"Your code is 12345"`
}
