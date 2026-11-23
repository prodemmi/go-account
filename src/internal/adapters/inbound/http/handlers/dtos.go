package handlers

// Request bodies for the public CAS API.

type sendOTPBody struct {
	Identifier string `json:"identifier" binding:"required" example:"+989130243519"`
}

type verifyOTPBody struct {
	Identifier string `json:"identifier" binding:"required" example:"+989130243519"`
	OTPCode    string `json:"otp_code" binding:"required" example:"12345"`
	OTPHash    string `json:"otp_hash" binding:"required" example:"fabecfb8ad5895a6a6309e5a7d1b5c48601d1f9e"`
}

type verifyOTTBody struct {
	OTT string `json:"ott" binding:"required" example:"9f2c..."`
}

type terminateSessionsBody struct {
	IDs []string `json:"ids" binding:"required" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
}
