package models

import uuid "github.com/satori/go.uuid"

type Token struct {
	BaseModel

	User   User
	UserID uuid.UUID
	App    App
	AppID  uuid.UUID
	Token  string
}
