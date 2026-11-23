package models

// App model info
// @Description App information
type App struct {
	BaseModel

	Name        string `json:"name" example:"Fotomi"`
	Url         string `json:"url" example:"https://fotomi.net/"`
	CallbackURL string `json:"callback_url" example:"https://fotomi.net/auth/callback"`
	WebhookURL  string `json:"webhook_url,omitempty" example:"https://fotomi.net/webhooks/cas"`

	OTPChannel string `json:"otp_channel" gorm:"default:sms" example:"sms"`
	OTPLength  int    `json:"otp_length" gorm:"default:5" example:"5"`
	OTPCharset string `json:"otp_charset" gorm:"default:numeric" example:"numeric"`
	OTPExpire  int    `json:"otp_expire" gorm:"default:120" example:"120"`
	OTTExpire  int    `json:"ott_expire" gorm:"default:60" example:"60"`

	SessionTimeout        int    `json:"session_timeout" gorm:"default:3600" example:"3600"`
	InfiniteSession       bool   `json:"infinite_session" gorm:"default:false" example:"false"`
	SessionCookieName     string `json:"session_cookie_name" gorm:"default:cas_cookie" example:"cas_cookie"`
	SessionCookieHTTPOnly bool   `json:"session_cookie_http_only" example:"true"`
	SessionCookieSecure   bool   `json:"session_cookie_secure" gorm:"default:false" example:"true"`
	SessionCookieSameSite string `json:"session_cookie_same_site" gorm:"default:lax" example:"lax"`
	RedirectBack          bool   `json:"redirect_back" gorm:"default:false" example:"true"`

	AuthPageEnabled             bool   `json:"auth_page_enabled" gorm:"default:false" example:"true"`
	AuthPageURL                 string `json:"auth_page_url" gorm:"default:null" example:"https://fotomi.net/auth/login"`
	AuthPageTitle               string `json:"auth_page_title" gorm:"default:null" example:"ورود به {app}"`
	AuthPageLayout              string `json:"auth_page_layout" gorm:"default:card" example:"card"`
	AuthPageFont                string `json:"auth_page_font" gorm:"default:vazirmatn" example:"vazirmatn"`
	AuthPageTextAlign           string `json:"auth_page_text_align" gorm:"default:center" example:"center"`
	AuthPageLogo                string `json:"auth_page_logo" gorm:"default:null" example:"https://fotomi.net/logo.png"`
	AuthPageLogoMask            string `json:"auth_page_logo_mask" gorm:"default:none" example:"circle"`
	AuthPageLogoPosition        string `json:"auth_page_logo_position" gorm:"default:top" example:"top"`
	AuthPagePrimaryColor        string `json:"auth_page_primary_color" gorm:"default:#3b82f6" example:"#3b82f6"`
	AuthPageBackgroundColor     string `json:"auth_page_background_color" gorm:"default:#f8fafc" example:"#f8fafc"`
	AuthPageCardBackgroundColor string `json:"auth_page_card_background_color" gorm:"default:#ffffff" example:"#ffffff"`
	AuthPagePinBackgroundColor  string `json:"auth_page_pin_background_color" gorm:"default:#f9fafb" example:"#f9fafb"`
	AuthPagePinTextColor        string `json:"auth_page_pin_text_color" gorm:"default:#111827" example:"#111827"`
	AuthPageTheme               string `json:"auth_page_theme" gorm:"default:light" example:"light"`
	AuthPageCaptchaEnabled      bool   `json:"auth_page_captcha_enabled" gorm:"default:false" example:"true"`
	AuthPageCaptchaSiteKey      string `json:"auth_page_captcha_site_key" gorm:"default:null" example:"public-site-key"`
	AuthPageCaptchaSecretKey    string `json:"auth_page_captcha_secret_key" gorm:"default:null"`
	AuthPageBannerMode          string `json:"auth_page_banner_mode" gorm:"default:none" example:"progress"`
	// AuthPageBanners stores the ordered banner list as JSON text.
	AuthPageBanners string `json:"auth_page_banners" gorm:"type:text" example:"[]"`
	// AuthPageCopyCustom means the copy fields below were explicitly edited and
	// empty strings must stay empty (hidden) instead of using defaults.
	AuthPageCopyCustom bool `json:"auth_page_copy_custom" gorm:"default:false"`

	AuthPageLoginTitle      string `json:"auth_page_login_title" gorm:"default:null" example:"ورود به {app}"`
	AuthPageLoginSubtitle   string `json:"auth_page_login_subtitle" gorm:"default:null" example:"کد ورود برای شما ارسال می‌شود."`
	AuthPageIdentifierLabel string `json:"auth_page_identifier_label" gorm:"default:null" example:"شماره موبایل"`
	AuthPageSendButton      string `json:"auth_page_send_button" gorm:"default:null" example:"دریافت کد ورود"`
	AuthPageVerifyTitle     string `json:"auth_page_verify_title" gorm:"default:null" example:"کد تایید را وارد کنید"`
	AuthPageVerifySubtitle  string `json:"auth_page_verify_subtitle" gorm:"default:null" example:"کد ارسال‌شده به {identifier} را وارد کنید."`
	AuthPageVerifyButton    string `json:"auth_page_verify_button" gorm:"default:null" example:"تایید و ورود"`
	AuthPageRetryText       string `json:"auth_page_retry_text" gorm:"default:null" example:"ارسال دوباره کد"`
	AuthPageBackText        string `json:"auth_page_back_text" gorm:"default:null" example:"تغییر ایمیل یا شماره موبایل"`
	AuthPageFooterText      string `json:"auth_page_footer_text" gorm:"default:null" example:"ورود امن با Aman CAS"`
}
