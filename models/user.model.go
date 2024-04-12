package models

import (
	uuid "github.com/satori/go.uuid"
)

// User model info
// @Description User information
type User struct {
	BaseModel

	// App id
	AppID uuid.UUID `json:"app_id,omitempty" gorm:"index" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`

	// Profile id
	ProfileID uuid.UUID `json:"profile_id,omitempty" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`

	// Current user ip
	IP string `json:"ip" example:"197.240.172.234"`

	// First user ip
	LoginIP string `json:"login_ip" example:"197.240.172.234"`

	// Stored API Token
	ApiToken *string `json:"api_token,omitempty" example:"eyJjbGl..."`
}
