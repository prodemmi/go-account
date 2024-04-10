package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Profile struct {
	BaseModel

	App  *App
	User *User

	AppID        uuid.UUID `json:"app_id"`
	UserID       uuid.UUID `json:"user_id"`
	Type         string    `json:"type"`
	Avatar       string    `json:"avatar"`
	AvatarStatus string    `json:"avatar_status"`
	Region       string    `json:"region"`
	Birthday     time.Time `json:"birthday"`
	Score        int       `json:"score"`
}
