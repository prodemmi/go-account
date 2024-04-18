package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Otp model info
// @Description Otp information
type Otp struct {
	BaseModel

	// App id
	AppID uuid.UUID `json:"app_id,omitempty" gorm:"index" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	// User id
	UserID uuid.UUID `json:"user_id,omitempty" example:"14a03569-d26b-4496-92e5-dfe8cb1855fe"`
	// code
	Code string `json:"code" example:"1234"`
	// url
	Type string `json:"type" gorm:"index" example:"email"`
	// Used
	Used bool `json:"used" gorm:"default:false" example:"true"`
	// Hash for security
	Hash *string `json:"hash" gorm:"default:null" example:"asnasd89asdasid"`
	// Expire at
	ExpireAt *time.Time `json:"expire_at" gorm:"default:null" example:"2012-05-22 07:20:02"`
}

func (otp *Otp) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("otp code: " + otp.Code)
	return
}
