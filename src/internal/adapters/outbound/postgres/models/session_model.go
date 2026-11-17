package models

import (
	"github.com/google/uuid"
	"time"
	"gorm.io/gorm"
)

// Session model info
// @Description Session information
type Session struct {
	BaseModel
	// App id
	AppID uuid.UUID `json:"app_id" gorm:"index" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	// session / service ticket token
	Token string `json:"token" gorm:"index;unique" example:"ticket-ST-..."`
	// expires at
	ExpiresAt time.Time `json:"expires_at" example:"2024-01-01T00:00:00Z"`
}

func (session *Session) AfterCreate(tx *gorm.DB) (err error) {
	return
}
