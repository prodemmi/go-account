package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Profile struct {
	BaseModel

	AppID  uuid.UUID `json:"app_id" gorm:"index"`
	UserID uuid.UUID `json:"user_id"`

	Type                           *string    `json:"type"`
	Avatar                         string     `json:"avatar"`
	AvatarStatus                   *string    `json:"avatar_status"`
	Province                       *string    `json:"province"`
	City                           *string    `json:"city"`
	Country                        *string    `json:"country" gorm:"index"`
	PostalCode                     *string    `json:"postal_code"`
	Timezone                       *string    `json:"timezone+"`
	Bio                            *string    `json:"bio"`
	Mobile                         *string    `json:"mobile"`
	FirstName                      string     `json:"first_name"`
	LastName                       string     `json:"last_name"`
	Email                          *string    `json:"email"`
	Password                       *string    `json:"-"`
	Address                        *string    `json:"address"`
	NationalCode                   *string    `json:"national_code"`
	IsActive                       *bool      `json:"is_active" gorm:"index;default:true"`
	Reagent                        *User      `json:"-"`
	ReagentID                      *uuid.UUID `json:"reagent_id" gorm:"foreignKey:User"`
	Gender                         *string    `json:"gender" gorm:"index"`
	Language                       *string    `json:"language" gorm:"index"`
	Birthday                       *time.Time `json:"birthday"`
	Score                          *uint      `json:"score" gorm:"default:0"`
	Theme                          *string    `json:"theme"`
	UserAgent                      *string    `json:"user_agent"`
	LoginCount                     *uint      `json:"login_count"`
	TwoFactorAuthenticationEnabled bool       `json:"two_factor_authentication_enabled"`
	EmailVerifiedAt                *time.Time `json:"email_verified_at"`
	LastActivity                   *time.Time `json:"last_activity"`
}
