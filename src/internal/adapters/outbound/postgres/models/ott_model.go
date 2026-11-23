package models

import (
	"time"

	"github.com/google/uuid"
)

// OTT model info
// @Description One-time ticket exchanged by an app for credentials
type OTT struct {
	BaseModel

	AppID      uuid.UUID `json:"app_id" gorm:"index"`
	SessionID  uuid.UUID `json:"session_id" gorm:"index"`
	Token      string    `json:"token" gorm:"index;unique"`
	Identifier string    `json:"identifier"`
	Used       bool      `json:"used" gorm:"default:false"`
	ExpireAt   time.Time `json:"expire_at" gorm:"index"`
}
