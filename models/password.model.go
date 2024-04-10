package models

import uuid "github.com/satori/go.uuid"

type Password struct {
	BaseModel

	App  *App
	User *User

	AppID    uuid.UUID `json:"app_id"`
	UserID   uuid.UUID `json:"user_id"`
	Password string    `json:"password"`
}
