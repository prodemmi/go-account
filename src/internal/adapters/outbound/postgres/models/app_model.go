package models

// App model info
// @Description App information
type App struct {
	BaseModel

	// name
	Name string `json:"name" example:"Fotomi"`
	// url
	Url string `json:"url" example:"https://fotomi.net/"`
	// url
	AuthUrl string `json:"auth_url" example:"https://fotomi.net/auth/callback"`
	// otp type
	OtpType string `json:"otp_type" gorm:"email" example:"email"`
	// session timeout in seconds
	SessionTimeout int `json:"session_timeout" gorm:"default:3600" example:"3600"`
	// custom auth page URL (if app hosts its own)
	AuthPageURL string `json:"auth_page_url" gorm:"default:null" example:"https://fotomi.net/auth/login"`
	// template name to render for branded auth pages
	AuthPageTemplate string `json:"auth_page_template" gorm:"default:default" example:"default"`
	// logo URL for branding
	AuthPageLogo string `json:"auth_page_logo" gorm:"default:null" example:"https://fotomi.net/logo.png"`
	// primary brand color (hex)
	AuthPagePrimaryColor string `json:"auth_page_primary_color" gorm:"default:#3b82f6" example:"#3b82f6"`
	// whether custom auth page is enabled
	AuthPageEnabled bool `json:"auth_page_enabled" gorm:"default:false" example:"true"`
}
