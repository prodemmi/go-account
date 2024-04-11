package models

import uuid "github.com/satori/go.uuid"

type User struct {
	BaseModel

	App   App       `json:"-"`
	AppID uuid.UUID `json:"app_id,omitempty"`

	Reagent   *User      `json:"-"`
	ReagentID *uuid.UUID `json:"reagent_id" gorm:"foreignKey:User"`

	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	Mobile       string  `json:"mobile"`
	NationalCode *string `json:"national_code"`
	IsActive     *bool   `json:"is_active" gorm:"index;default:true"`
	SignupIP     string  `json:"signup_ip"`
	LoginIP      string  `json:"login_ip"`
	ApiToken     *string `json:"api_token,omitempty"`
	Password     *string `json:"-"`
}
