package app

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// OTP delivery channels an app can be configured with.
const (
	OTPChannelSMS   = "sms"
	OTPChannelEmail = "email"
	OTPChannelBoth  = "both"
)

// Character sets an app can use for its OTP codes.
const (
	OTPCharsetNumeric      = "numeric"
	OTPCharsetAlpha        = "alpha"
	OTPCharsetAlphanumeric = "alphanumeric"
)

// Auth page themes.
const (
	ThemeLight = "light"
	ThemeDark  = "dark"
)

// Auth page layouts.
const (
	LayoutCard    = "card"
	LayoutSplit   = "split"
	LayoutMinimal = "minimal"
)

// Auth page font choices. The Persian families are bundled with the binary,
// the Latin ones cover English copy and identifiers.
const (
	FontVazirmatn = "vazirmatn"
	FontIRANSans  = "iransans"
	FontEstedad   = "estedad"
	FontSahel     = "sahel"
	FontInter     = "inter"
	FontManrope   = "manrope"
	FontSystem    = "system"
)

// Text alignment choices for the auth page.
const (
	AlignStart  = "start"
	AlignCenter = "center"
	AlignEnd    = "end"
)

// Logo mask choices.
const (
	LogoMaskNone    = "none"
	LogoMaskCircle  = "circle"
	LogoMaskRounded = "rounded"
	LogoMaskSquare  = "square"
)

// Logo position choices. Each one maps to a distinct place in the auth card:
// above the card, beside the title (start/end) or at the bottom of the card.
const (
	LogoPositionTop    = "top"
	LogoPositionStart  = "start"
	LogoPositionEnd    = "end"
	LogoPositionCenter = "center"
	LogoPositionBottom = "bottom"
)

// Cookie SameSite choices.
const (
	SameSiteDefault = "default"
	SameSiteLax     = "lax"
	SameSiteStrict  = "strict"
	SameSiteNone    = "none"
)

// Banner duration display modes. They decide how the remaining time of the
// current banner is shown while it is visible.
const (
	BannerModeNone     = "none"
	BannerModeProgress = "progress"
	BannerModeTime     = "time"
)

// DefaultSessionCookieName is used when an app does not configure a cookie name.
const DefaultSessionCookieName = "cas_cookie"

// InfiniteSessionTTL is the session lifetime used when an app marks its
// sessions as infinite. CAS does not manage users, so "infinite" simply means
// a very long expiry that is never reached in practice.
const InfiniteSessionTTL = 100 * 365 * 24 * time.Hour

// Banner is one auth-page banner. Banners are ordered by Order and each one is
// shown for Duration seconds before the next one slides in.
type Banner struct {
	URL      string `json:"url"`
	Text     string `json:"text"`
	Link     string `json:"link,omitempty"`
	Duration int    `json:"duration"`
	Order    int    `json:"order"`
}

// App is the domain entity representing a client application in the CAS.
// CAS only authenticates users; each app configures the whole auth experience.
type App struct {
	ID          uuid.UUID
	Name        string
	Url         string
	CallbackURL string // where CAS redirects the browser after a successful login
	WebhookURL  string // receives an event after every auth action

	// OTP configuration
	OTPChannel string        // sms | email | both
	OTPLength  int           // 4..32
	OTPCharset string        // numeric | alpha | alphanumeric
	OTPExpire  time.Duration // how long an OTP code is valid
	OTTExpire  time.Duration // how long a one-time ticket is valid

	// Session configuration
	SessionTimeout        time.Duration
	InfiniteSession       bool
	SessionCookieName     string // default cas_cookie
	SessionCookieHTTPOnly bool
	SessionCookieSecure   bool
	SessionCookieSameSite string // default | lax | strict | none
	RedirectBack          bool   // return to the CAS entry URL when it shares the app origin

	// Auth page customization (CAS serves branded auth pages per app)
	AuthPageEnabled             bool
	AuthPageURL                 string // custom auth page URL (if app hosts its own)
	AuthPageTitle               string // browser document title
	AuthPageLayout              string // card | split | minimal
	AuthPageFont                string // vazirmatn | iransans | system
	AuthPageTextAlign           string // start | center | end
	AuthPageLogo                string // logo URL for branding
	AuthPageLogoMask            string // none | circle | rounded | square
	AuthPageLogoPosition        string // top | start | end | center
	AuthPagePrimaryColor        string // primary brand color (hex)
	AuthPageBackgroundColor     string // page background color (hex)
	AuthPageCardBackgroundColor string // card/panel background color (hex)
	AuthPagePinBackgroundColor  string // OTP box background color (hex)
	AuthPagePinTextColor        string // OTP box text color (hex)
	AuthPageTheme               string // light | dark
	AuthPageCaptchaEnabled      bool
	AuthPageCaptchaSiteKey      string // public ArCaptcha site key
	AuthPageCaptchaSecretKey    string // server-side ArCaptcha secret
	AuthPageBannerMode          string // none | progress | time
	AuthPageBanners             []Banner

	// Custom auth page copy. Empty values fall back to the defaults below and
	// "{app}" expands to the app name. AuthPageCopyCustom is set once the
	// admin edits the copy section: from then on an empty field means "hide
	// this text" instead of "use the default".
	AuthPageCopyCustom      bool
	AuthPageLoginTitle      string
	AuthPageLoginSubtitle   string
	AuthPageIdentifierLabel string
	AuthPageSendButton      string
	AuthPageVerifyTitle     string
	AuthPageVerifySubtitle  string
	AuthPageVerifyButton    string
	AuthPageRetryText       string
	AuthPageBackText        string
	AuthPageFooterText      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// AllowsMobile reports whether the app accepts a mobile number as identifier.
func (a App) AllowsMobile() bool {
	return a.OTPChannel == "" || a.OTPChannel == OTPChannelSMS || a.OTPChannel == OTPChannelBoth
}

// AllowsEmail reports whether the app accepts an email address as identifier.
func (a App) AllowsEmail() bool {
	return a.OTPChannel == OTPChannelEmail || a.OTPChannel == OTPChannelBoth
}

// OTPInputType returns the HTML input type for the identifier field. It is
// strict: a mobile-only app gets type=tel, an email-only app type=email and a
// mixed app a plain text field.
func (a App) OTPInputType() string {
	switch {
	case a.AllowsEmail() && a.AllowsMobile():
		return "text"
	case a.AllowsEmail():
		return "email"
	default:
		return "tel"
	}
}

// OTPCodeInputType returns the HTML type for an OTP character box.
func (a App) OTPCodeInputType() string {
	if a.OTPCharset == OTPCharsetNumeric {
		return "tel"
	}
	return "text"
}

// SessionTTL returns the lifetime used for newly created sessions.
func (a App) SessionTTL() time.Duration {
	if a.InfiniteSession {
		return InfiniteSessionTTL
	}
	if a.SessionTimeout <= 0 {
		return time.Hour
	}
	return a.SessionTimeout
}

// ExpandText replaces the supported placeholders in a customized string.
func (a App) ExpandText(s string) string {
	s = strings.ReplaceAll(s, "{app}", a.Name)
	s = strings.ReplaceAll(s, "{name}", a.Name)
	return s
}

// Text returns a customized auth-page string, falling back to the default and
// expanding placeholders.
func (a App) Text(custom, fallback string) string {
	if strings.TrimSpace(custom) == "" {
		custom = fallback
	}
	return a.ExpandText(custom)
}

// CopyText is Text for the copy editor. Once the admin has edited the copy
// section, an empty value is intentional and the string is hidden instead of
// falling back to the built-in default.
func (a App) CopyText(custom, fallback string) string {
	if !a.AuthPageCopyCustom && strings.TrimSpace(custom) == "" {
		custom = fallback
	}
	return a.ExpandText(custom)
}

// Normalize fills in safe defaults for values an app may not have configured.
// It is applied before an app is persisted and mirrored when an app is loaded.
func (a App) Normalize() App {
	a.Name = strings.TrimSpace(a.Name)
	a.Url = strings.TrimSpace(a.Url)
	a.CallbackURL = strings.TrimSpace(a.CallbackURL)
	a.WebhookURL = strings.TrimSpace(a.WebhookURL)

	switch a.OTPChannel {
	case OTPChannelSMS, OTPChannelEmail, OTPChannelBoth:
	default:
		a.OTPChannel = OTPChannelSMS
	}

	if a.OTPLength < 4 || a.OTPLength > 32 {
		a.OTPLength = 5
	}
	switch a.OTPCharset {
	case OTPCharsetNumeric, OTPCharsetAlpha, OTPCharsetAlphanumeric:
	default:
		a.OTPCharset = OTPCharsetNumeric
	}

	if a.OTPExpire <= 0 {
		a.OTPExpire = 2 * time.Minute
	}
	if a.OTTExpire <= 0 {
		a.OTTExpire = time.Minute
	}
	if a.SessionTimeout <= 0 {
		a.SessionTimeout = time.Hour
	}

	if a.SessionCookieName = strings.TrimSpace(a.SessionCookieName); a.SessionCookieName == "" {
		a.SessionCookieName = DefaultSessionCookieName
		a.SessionCookieHTTPOnly = true
	}
	switch a.SessionCookieSameSite {
	case SameSiteLax, SameSiteStrict, SameSiteNone, SameSiteDefault:
	default:
		a.SessionCookieSameSite = SameSiteLax
	}

	if a.AuthPageLayout == "" {
		a.AuthPageLayout = LayoutCard
	}
	switch a.AuthPageLayout {
	case LayoutCard, LayoutSplit, LayoutMinimal:
	default:
		a.AuthPageLayout = LayoutCard
	}
	if a.AuthPageFont == "" {
		a.AuthPageFont = FontVazirmatn
	}
	switch a.AuthPageFont {
	case FontVazirmatn, FontIRANSans, FontEstedad, FontSahel, FontInter, FontManrope, FontSystem:
	default:
		a.AuthPageFont = FontVazirmatn
	}
	if a.AuthPageTextAlign == "" {
		a.AuthPageTextAlign = AlignCenter
	}
	switch a.AuthPageTextAlign {
	case AlignStart, AlignCenter, AlignEnd:
	default:
		a.AuthPageTextAlign = AlignCenter
	}
	if a.AuthPageLogoMask == "" {
		a.AuthPageLogoMask = LogoMaskNone
	}
	switch a.AuthPageLogoMask {
	case LogoMaskNone, LogoMaskCircle, LogoMaskRounded, LogoMaskSquare:
	default:
		a.AuthPageLogoMask = LogoMaskNone
	}
	if a.AuthPageLogoPosition == "" {
		a.AuthPageLogoPosition = LogoPositionTop
	}
	switch a.AuthPageLogoPosition {
	case LogoPositionTop, LogoPositionStart, LogoPositionEnd, LogoPositionCenter, LogoPositionBottom:
	default:
		a.AuthPageLogoPosition = LogoPositionTop
	}
	if a.AuthPageTheme != ThemeDark {
		a.AuthPageTheme = ThemeLight
	}
	if a.AuthPageTitle == "" && !a.AuthPageCopyCustom {
		a.AuthPageTitle = "ورود به {app}"
	}
	if a.AuthPagePrimaryColor == "" {
		a.AuthPagePrimaryColor = "#3b82f6"
	}
	if a.AuthPageBackgroundColor == "" {
		if a.AuthPageTheme == ThemeDark {
			a.AuthPageBackgroundColor = "#0f172a"
		} else {
			a.AuthPageBackgroundColor = "#f8fafc"
		}
	}
	if a.AuthPageCardBackgroundColor == "" {
		if a.AuthPageTheme == ThemeDark {
			a.AuthPageCardBackgroundColor = "#0f172a"
		} else {
			a.AuthPageCardBackgroundColor = "#ffffff"
		}
	}
	if a.AuthPagePinBackgroundColor == "" {
		if a.AuthPageTheme == ThemeDark {
			a.AuthPagePinBackgroundColor = "#1e293b"
		} else {
			a.AuthPagePinBackgroundColor = "#f9fafb"
		}
	}
	if a.AuthPagePinTextColor == "" {
		if a.AuthPageTheme == ThemeDark {
			a.AuthPagePinTextColor = "#e2e8f0"
		} else {
			a.AuthPagePinTextColor = "#111827"
		}
	}
	switch a.AuthPageBannerMode {
	case BannerModeNone, BannerModeProgress, BannerModeTime:
	default:
		a.AuthPageBannerMode = BannerModeNone
	}

	for i := range a.AuthPageBanners {
		a.AuthPageBanners[i].URL = strings.TrimSpace(a.AuthPageBanners[i].URL)
		a.AuthPageBanners[i].Text = strings.TrimSpace(a.AuthPageBanners[i].Text)
		a.AuthPageBanners[i].Link = strings.TrimSpace(a.AuthPageBanners[i].Link)
		if a.AuthPageBanners[i].Duration <= 0 {
			a.AuthPageBanners[i].Duration = 6
		}
	}

	return a
}
