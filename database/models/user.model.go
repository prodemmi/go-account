package models

import (
	"github.com/google/uuid"
	"time"
)

// User model info
// @Description User information
type User struct {
	BaseModel

	// App id
	AppID uuid.UUID `json:"app_id,omitempty" gorm:"index" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`

	// Current user ip
	IP string `json:"ip" example:"197.240.172.234"`

	// First user ip
	LoginIP string `json:"login_ip" example:"197.240.172.234"`

	// User type for example admin, guest, vip ...
	Type *string `json:"type" example:"admin"`

	// User avatar url
	Avatar                         string     `json:"avatar" example:"https://avatars.githubusercontent.com/u/55631515?v=4"`
	AvatarStatus                   *string    `json:"avatar_status" example:"at work"`
	Province                       *string    `json:"province" example:"Isfahan"`
	City                           *string    `json:"city" example:"Isfahan"`
	Country                        *string    `json:"country" gorm:"index" example:"Iran"`
	PostalCode                     *string    `json:"postal_code" example:"8457954212"`
	Timezone                       *string    `json:"timezone+" example:"Indian/Antananarivo"`
	Bio                            *string    `json:"bio" example:"Hello my name is emad."`
	Mobile                         *string    `json:"mobile" gorm:"unique" example:"+989130243519"`
	FirstName                      string     `json:"first_name" example:"Emad"`
	LastName                       string     `json:"last_name" example:"Malekpour"`
	Email                          *string    `json:"email" gorm:"unique" example:"prodemmi@gmail.com"`
	Password                       string     `json:"-"`
	Address                        *string    `json:"address" example:"417 Kennedy Isle Apt. 706"`
	NationalCode                   *string    `json:"national_code" example:"1273039045"`
	IsActive                       *bool      `json:"is_active" gorm:"index;default:true" example:"true"`
	Reagent                        *User      `json:"-"`
	ReagentID                      *uuid.UUID `json:"reagent_id" gorm:"foreignKey:User" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	Gender                         *string    `json:"gender" gorm:"index" example:"Male"`
	Language                       *string    `json:"language" gorm:"index" example:"fa"`
	Birthday                       *time.Time `json:"birthday" example:"2012-05-22 07:20:02"`
	Score                          *uint      `json:"score" gorm:"default:0" example:"12"`
	Theme                          *string    `json:"theme" example:"dracula"`
	UserAgent                      *string    `json:"user_agent" example:"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_5_6) AppleWebKit/5311 (KHTML, like Gecko) Chrome/40.0.877.0 Safari/5311"`
	LoginCount                     *uint      `json:"login_count" example:"24"`
	TwoFactorAuthenticationEnabled bool       `json:"two_factor_authentication_enabled" example:"false"`
	EmailVerifiedAt                *time.Time `json:"email_verified_at" example:"1986-01-24 00:10:19"`
	LastActivity                   *time.Time `json:"last_activity" example:"1986-01-24 00:10:19"`

	SignedUp bool `json:"-" gorm:"default:false" example:"false"`
	// Stored API Token
	ApiToken *string `json:"-" example:"eyJjbGl..."`
}
