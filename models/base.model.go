package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	CreatedAt time.Time  `json:"created_at" example:"1973-11-06 14:12:31"`
	UpdatedAt time.Time  `json:"updated_at" example:"1973-11-06 14:12:31"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at" example:"null"`
}
