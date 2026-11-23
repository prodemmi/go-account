package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Otp model info
// @Description Otp information
type Otp struct {
	BaseModel

	AppID      uuid.UUID  `json:"app_id,omitempty" gorm:"index" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	Identifier string     `json:"identifier" gorm:"index" example:"+989130243519"`
	Code       string     `json:"code" example:"1234"`
	Type       string     `json:"type" gorm:"index" example:"email"`
	Used       bool       `json:"used" gorm:"default:false" example:"true"`
	Hash       *string    `json:"hash" gorm:"default:null" example:"asnasd89asdasid"`
	ExpireAt   *time.Time `json:"expire_at" gorm:"default:null" example:"2012-05-22 07:20:02"`
}

func (otp *Otp) AfterCreate(_ *gorm.DB) (err error) {
	fmt.Println("otp code: " + otp.Code)
	return
}
