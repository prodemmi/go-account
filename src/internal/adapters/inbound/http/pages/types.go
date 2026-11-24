package pages

import (
	"io"

	"aman/internal/domains/app"
)

// AuthData holds the data passed to auth page components.
type AuthData struct {
	App                 app.App
	AppID               string
	Identifier          string
	Error               string
	Success             string
	OtpHash             string
	CsrfToken           string
	ReturnTo            string
	PageTitle           string
	DocumentTitle       string
	PrimaryColor        string
	BackgroundColor     string
	CardBackgroundColor string
	TextColor           string
	MutedTextColor      string
	PinBackgroundColor  string
	PinTextColor        string
	Theme               string
	Font                string
	TextAlign           string
	Layout              string
	LogoURL             string
	LogoMask            string
	LogoPosition        string
	BannerMode          string
	Banners             []app.Banner
	CopyCustom          bool
	AllowMobile         bool
	AllowEmail          bool
	Step                string

	Title           string
	Subtitle        string
	IdentifierLabel string
	SendButton      string
	VerifyTitle     string
	VerifySubtitle  string
	VerifyButton    string
	RetryText       string
	BackText        string
	FooterText      string

	CaptchaEnabled  bool
	CaptchaProvider string
	CaptchaSiteKey  string

	OTPLength        int
	OTPInputType     string
	OTPCodeInputType string
	OTPExpire        int
	OTPExpiresAt     int64
}

func EditIdentifierText(allowEmail, allowMobile bool) string {
	switch {
	case allowEmail && allowMobile:
		return "تغییر ایمیل یا شماره موبایل"
	case allowEmail:
		return "تغییر ایمیل"
	default:
		return "تغییر شماره موبایل"
	}
}

// WriteEscaped writes HTML-escaped text to w.
func WriteEscaped(w io.Writer, s string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			w.Write([]byte("&amp;"))
		case '<':
			w.Write([]byte("&lt;"))
		case '>':
			w.Write([]byte("&gt;"))
		case '"':
			w.Write([]byte("&quot;"))
		case '\'':
			w.Write([]byte("&#39;"))
		default:
			w.Write([]byte{s[i]})
		}
	}
}
