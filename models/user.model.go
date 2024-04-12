package models

import (
	uuid "github.com/satori/go.uuid"
)

type User struct {
	BaseModel

	AppID     uuid.UUID `json:"app_id,omitempty" gorm:"index"`
	ProfileID uuid.UUID `json:"profile_id,omitempty"`
	IP        string    `json:"ip"`
	LoginIP   string    `json:"login_ip"`
	ApiToken  *string   `json:"api_token,omitempty"`
}
