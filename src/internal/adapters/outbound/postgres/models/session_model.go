package models

import (
	"time"

	"github.com/google/uuid"
)

// Session model info
// @Description Session information
type Session struct {
	BaseModel

	AppID uuid.UUID `json:"app_id" gorm:"index" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	// session / service ticket token
	Token string `json:"token" gorm:"index;unique" example:"ticket-ST-..."`
	// mobile or email used to authenticate
	Identifier string `json:"identifier" gorm:"index" example:"+989130243519"`
	Mobile     string `json:"mobile,omitempty" example:"+989130243519"`
	Email      string `json:"email,omitempty" example:"user@example.com"`
	UserAgent  string `json:"user_agent,omitempty"`
	IP         string `json:"ip,omitempty" gorm:"index"`
	Country    string `json:"country,omitempty"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	Device     string `json:"device,omitempty"`
	// pending | active | revoked. Expired is derived from ExpiresAt.
	// No DB default: old rows stay empty and are read as active/revoked from
	// the legacy is_logged_in flag.
	Status     string `json:"status" gorm:"index" example:"active"`
	IsLoggedIn bool   `json:"is_logged_in" gorm:"index;default:false"`
	Infinite   bool   `json:"infinite" gorm:"default:false"`

	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	ExpiresAt   time.Time  `json:"expires_at" example:"2024-01-01T00:00:00Z"`
}
