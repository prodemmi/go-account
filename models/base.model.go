package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
