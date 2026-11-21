package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"aman/internal/adapters/inbound/http/pages"
	auth "aman/internal/adapters/inbound/http/pages/auth"
	"aman/internal/domains/app"
	"aman/internal/domains/ott"
	"aman/internal/domains/session"
	"aman/internal/ports/inbound"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

const returnToCookie = "cas_return_to"

// AuthPageHandler serves the branded auth pages (CAS responsibility).
// CAS is the central authentication service: apps redirect users here for
// auth, and CAS renders app-branded login/verify-OTP pages.
type AuthPageHandler struct {
	apps inbound.AppUseCase
	auth inbound.AuthUseCase
}

func NewAuthPageHandler(apps inbound.AppUseCase, auth inbound.AuthUseCase) *AuthPageHandler {
	return &AuthPageHandler{apps: apps, auth: auth}
}

// Login renders the branded login page for an app.
func (h AuthPageHandler) Login(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		redirectNotFound(c)
		return
	}

	a, err := h.apps.Get(c.Request.Context(), appID)
	if err != nil {
		redirectNotFound(c)
		return
	}

	// If the app hosts its own auth page, send the user there.
	if a.AuthPageEnabled && a.AuthPageURL != "" && isSafeRedirectURL(a.AuthPageURL) {
		target := a.AuthPageURL
		separator := "?"
		if strings.Contains(target, "?") {
			separator = "&"
		}
		c.Redirect(http.StatusTemporaryRedirect, target+separator+"app_id="+url.QueryEscape(appID))
		return
	}

	returnTo := h.rememberReturnTo(c, a)
	data := authData(*a, "ورود", c.Query("error"))
	data.Step = "login"
	data.ReturnTo = returnTo
	if a.AuthPageCaptchaEnabled && arCaptchaConfigured(a) {
		data.CaptchaEnabled = true
		data.CaptchaSiteKey = a.AuthPageCaptchaSiteKey
	}

	renderPage(c, auth.Login(data))
}

// VerifyOTP renders the branded code verification page.
func (h AuthPageHandler) VerifyOTPPage(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		redirectNotFound(c)
		return
	}

	a, err := h.apps.Get(c.Request.Context(), appID)
	if err != nil {
		redirectNotFound(c)
		return
	}

	data := authData(*a, "تایید هویت", c.Query("error"))
	data.Step = "verify"
	data.Identifier = c.Query("identifier")
	data.OtpHash = c.Query("otp_hash")
	data.ReturnTo = h.resolveReturnTo(c, a)
	if a.AuthPageCaptchaEnabled && arCaptchaConfigured(a) {
		data.CaptchaEnabled = true
		data.CaptchaSiteKey = a.AuthPageCaptchaSiteKey
	}
	if expires := c.Query("expires_at"); expires != "" {
		if unix, err := strconv.ParseInt(expires, 10, 64); err == nil {
			data.OTPExpiresAt = unix
		}
	}
	renderPage(c, auth.VerifyOTP(data))
}

// SendOTP handles the login form: it creates and sends the OTP, then redirects
// the browser to the verification page.
func (h AuthPageHandler) SendOTP(c *gin.Context) {
	appID := c.PostForm("app_id")
	identifier := strings.TrimSpace(c.PostForm("identifier"))
	if appID == "" {
		redirectNotFound(c)
		return
	}

	a, err := h.apps.Get(c.Request.Context(), appID)
	if err != nil {
		redirectNotFound(c)
		return
	}

	returnTo := h.rememberReturnToValue(c, a, c.PostForm("redirect"))
	loginURL := loginPageURLWithReturn(a, appID, returnTo)

	if a.AuthPageCaptchaEnabled {
		// The retry button on the verify step is only enabled after the code
		// countdown ends, so it is exempt from a second captcha challenge.
		if c.PostForm("retry") != "1" && arCaptchaConfigured(a) {
			if !verifyArCaptcha(a, c.PostForm("arcaptcha-token")) {
				authErrorReply(c, loginURL, "تایید کپچا ناموفق بود. دوباره تلاش کنید.")
				return
			}
		}
	}

	meta := metaFrom(c)
	challenge, err := h.auth.SendOTP(c.Request.Context(), inbound.SendOTPCommand{
		AppID:      appID,
		Identifier: identifier,
		UserAgent:  meta.UserAgent,
		IP:         meta.IP,
		Country:    meta.Country,
		City:       meta.City,
		Region:     meta.Region,
		Device:     meta.Device,
	})
	if err != nil {
		authErrorReply(c, loginURL, humanizeAuthError(err))
		return
	}

	if wantsSPA(c) {
		c.JSON(http.StatusOK, gin.H{
			"ok":         true,
			"step":       "verify",
			"identifier": challenge.Identifier,
			"otp_hash":   challenge.OTPHash,
			"length":     challenge.Length,
			"expires_at": challenge.ExpiresAt.Unix(),
			"redirect":   returnTo,
		})
		return
	}

	target := "/auth/verify-otp?app_id=" + appID +
		"&identifier=" + queryEscape(challenge.Identifier) +
		"&otp_hash=" + queryEscape(challenge.OTPHash) +
		"&expires_at=" + strconv.FormatInt(challenge.ExpiresAt.Unix(), 10)
	if returnTo != "" {
		target += "&redirect=" + queryEscape(returnTo)
	}
	c.Redirect(http.StatusFound, target)
}

// VerifyOTP handles the verification form. On success it sets the CAS cookie
// and redirects to the app callback (or the remembered return URL) with ?ott=.
func (h AuthPageHandler) VerifyOTP(c *gin.Context) {
	appID := c.PostForm("app_id")
	identifier := strings.TrimSpace(c.PostForm("identifier"))
	otpHash := c.PostForm("otp_hash")
	code := strings.TrimSpace(c.PostForm("otp_code"))

	if appID == "" {
		redirectNotFound(c)
		return
	}

	a, err := h.apps.Get(c.Request.Context(), appID)
	if err != nil {
		redirectNotFound(c)
		return
	}

	returnTo := h.rememberReturnToValue(c, a, c.PostForm("redirect"))
	backURL := "/auth/verify-otp?app_id=" + appID +
		"&identifier=" + queryEscape(identifier) +
		"&otp_hash=" + queryEscape(otpHash)
	if expires := c.PostForm("expires_at"); expires != "" {
		if unix, err := strconv.ParseInt(expires, 10, 64); err == nil {
			backURL += "&expires_at=" + strconv.FormatInt(unix, 10)
		}
	}
	if returnTo != "" {
		backURL += "&redirect=" + queryEscape(returnTo)
	}

	meta := metaFrom(c)
	result, err := h.auth.VerifyOTP(c.Request.Context(), inbound.VerifyOTPCommand{
		AppID:      appID,
		Identifier: identifier,
		Code:       code,
		OTPHash:    otpHash,
		UserAgent:  meta.UserAgent,
		IP:         meta.IP,
		Country:    meta.Country,
		City:       meta.City,
		Region:     meta.Region,
		Device:     meta.Device,
		ReturnTo:   returnTo,
	})
	if err != nil {
		if errors.Is(err, session.ErrForbidden) || errors.Is(err, ott.ErrForbidden) || errors.Is(err, app.ErrForbidden) {
			if wantsSPA(c) {
				c.JSON(http.StatusOK, gin.H{"ok": false, "error": "این نشست معتبر نیست."})
				return
			}
			redirectNotFound(c)
			return
		}
		authErrorReply(c, backURL, humanizeAuthError(err))
		return
	}

	h.clearReturnTo(c)
	setSessionCookie(c, result.App, result.Session.Token, result.Session.ExpiresAt, result.Session.Infinite)

	if wantsSPA(c) {
		c.JSON(http.StatusOK, gin.H{"ok": true, "redirect": result.RedirectURL})
		return
	}
	c.Redirect(http.StatusFound, result.RedirectURL)
}

// Check is the entry point apps send the browser to. It reuses a valid CAS
// session by issuing a fresh OTT; otherwise it redirects to the auth page.
func (h AuthPageHandler) Check(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		redirectNotFound(c)
		return
	}

	a, err := h.apps.Get(c.Request.Context(), appID)
	if err != nil {
		redirectNotFound(c)
		return
	}

	returnTo := h.rememberReturnTo(c, a)
	result, err := h.auth.Check(c.Request.Context(), inbound.CheckCommand{
		AppID:        appID,
		SessionToken: sessionCookie(c, a),
		UserAgent:    c.Request.UserAgent(),
		ReturnTo:     returnTo,
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	if result.Authenticated {
		c.Redirect(http.StatusFound, result.RedirectURL)
		return
	}

	c.Redirect(http.StatusFound, loginPageURLWithReturn(a, appID, returnTo))
}

// Logout revokes the CAS session and sends the browser back to the app.
func (h AuthPageHandler) Logout(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		clearSessionCookies(c, &app.App{})
		c.Redirect(http.StatusFound, "/")
		return
	}

	a, err := h.apps.Get(c.Request.Context(), appID)
	if err != nil {
		redirectNotFound(c)
		return
	}

	if err := h.auth.Logout(c.Request.Context(), inbound.LogoutCommand{
		AppID:        appID,
		SessionToken: sessionCookie(c, a),
	}); err != nil {
		redirectNotFound(c)
		return
	}

	clearSessionCookies(c, a)
	c.Redirect(http.StatusFound, a.Url)
}

func authData(a app.App, title, errMessage string) pages.AuthData {
	a = a.Normalize()

	background := a.AuthPageBackgroundColor
	if background == "" {
		if a.AuthPageTheme == app.ThemeDark {
			background = "#0f172a"
		} else {
			background = "#f8fafc"
		}
	}

	return pages.AuthData{
		App:                 a,
		AppID:               a.ID.String(),
		PageTitle:           title,
		DocumentTitle:       a.Text(a.AuthPageTitle, "ورود به {app}"),
		Error:               errMessage,
		PrimaryColor:        a.AuthPagePrimaryColor,
		BackgroundColor:     background,
		CardBackgroundColor: a.AuthPageCardBackgroundColor,
		PinBackgroundColor:  a.AuthPagePinBackgroundColor,
		PinTextColor:        a.AuthPagePinTextColor,
		Theme:               a.AuthPageTheme,
		Font:                a.AuthPageFont,
		TextAlign:           a.AuthPageTextAlign,
		Layout:              a.AuthPageLayout,
		LogoURL:             a.AuthPageLogo,
		LogoMask:            a.AuthPageLogoMask,
		LogoPosition:        a.AuthPageLogoPosition,
		BannerMode:          a.AuthPageBannerMode,
		Banners:             a.AuthPageBanners,
		CopyCustom:          a.AuthPageCopyCustom,
		AllowMobile:         a.AllowsMobile(),
		AllowEmail:          a.AllowsEmail(),
		Title:               a.CopyText(a.AuthPageLoginTitle, "ورود به {app}"),
		Subtitle:            a.CopyText(a.AuthPageLoginSubtitle, ""),
		IdentifierLabel:     a.CopyText(a.AuthPageIdentifierLabel, ""),
		SendButton:          a.CopyText(a.AuthPageSendButton, "دریافت کد ورود"),
		VerifyTitle:         a.CopyText(a.AuthPageVerifyTitle, "کد تایید را وارد کنید"),
		VerifySubtitle:      a.CopyText(a.AuthPageVerifySubtitle, ""),
		VerifyButton:        a.CopyText(a.AuthPageVerifyButton, "تایید و ورود"),
		RetryText:           a.CopyText(a.AuthPageRetryText, "ارسال دوباره کد"),
		BackText:            a.CopyText(a.AuthPageBackText, "تغییر ایمیل یا شماره موبایل"),
		FooterText:          a.CopyText(a.AuthPageFooterText, "ورود امن با Aman CAS"),
		OTPLength:           a.OTPLength,
		OTPInputType:        a.OTPInputType(),
		OTPCodeInputType:    a.OTPCodeInputType(),
		OTPExpire:           int(a.OTPExpire.Seconds()),
	}
}

func renderPage(c *gin.Context, component templ.Component) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		redirectNotFound(c)
	}
}

// rememberReturnTo stores the URL the user came from when it shares the app's
// origin. Without this, a successful login always lands on the callback URL.
func (h AuthPageHandler) rememberReturnTo(c *gin.Context, a *app.App) string {
	candidate := strings.TrimSpace(c.Query("redirect"))
	if candidate == "" {
		candidate = strings.TrimSpace(c.GetHeader("Referer"))
	}
	return h.rememberReturnToValue(c, a, candidate)
}

func (h AuthPageHandler) rememberReturnToValue(c *gin.Context, a *app.App, candidate string) string {
	if !a.RedirectBack {
		return ""
	}
	if valid := safeReturnTo(a, candidate); valid != "" {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     returnToCookie,
			Value:    url.QueryEscape(valid),
			Path:     "/",
			MaxAge:   int((30 * time.Minute).Seconds()),
			HttpOnly: true,
			Secure:   isSecureRequest(c),
			SameSite: http.SameSiteLaxMode,
		})
		return valid
	}
	return h.resolveReturnTo(c, a)
}

func (h AuthPageHandler) resolveReturnTo(c *gin.Context, a *app.App) string {
	if !a.RedirectBack {
		return ""
	}
	raw, err := c.Cookie(returnToCookie)
	if err != nil || raw == "" {
		return ""
	}
	decoded, err := url.QueryUnescape(raw)
	if err != nil {
		return ""
	}
	return safeReturnTo(a, decoded)
}

func (h AuthPageHandler) clearReturnTo(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     returnToCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   isSecureRequest(c),
		SameSite: http.SameSiteLaxMode,
	})
}

// safeReturnTo only accepts absolute URLs on the app's own origin, so the CAS
// cannot be used as an open redirect.
func safeReturnTo(a *app.App, candidate string) string {
	if candidate == "" {
		return ""
	}
	appURL, err := url.Parse(strings.TrimSpace(a.Url))
	if err != nil || appURL.Scheme == "" || appURL.Host == "" {
		return ""
	}
	target, err := url.Parse(candidate)
	if err != nil || target.Scheme == "" || target.Host == "" {
		return ""
	}
	if !strings.EqualFold(appURL.Scheme, target.Scheme) || !strings.EqualFold(appURL.Host, target.Host) {
		return ""
	}
	return target.String()
}

func loginPageURLWithReturn(a *app.App, appID, returnTo string) string {
	target := "/auth/login?app_id=" + appID
	if returnTo != "" {
		target += "&redirect=" + queryEscape(returnTo)
	}
	return target
}

func humanizeAuthError(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, session.ErrForbidden), errors.Is(err, ott.ErrForbidden):
		return "این نشست معتبر نیست."
	case strings.Contains(err.Error(), "mobile"):
		return "این برنامه فقط شماره موبایل را می‌پذیرد."
	case strings.Contains(err.Error(), "email"):
		return "این برنامه فقط ایمیل را می‌پذیرد."
	case strings.Contains(err.Error(), "otp"):
		return "کد وارد شده معتبر نیست یا منقضی شده است."
	default:
		return "خطایی رخ داد، دوباره تلاش کنید."
	}
}

// wantsSPA reports whether the browser asked for a JSON response so the auth
// page can advance between steps without a full reload. Without JavaScript the
// same handlers keep working through normal form posts and redirects.
func wantsSPA(c *gin.Context) bool {
	return c.GetHeader("X-Aman-SPA") == "1" ||
		strings.Contains(c.GetHeader("Accept"), "application/json")
}

// authErrorReply returns the error in the format the caller can consume: JSON
// for the SPA, a redirect with ?error= for the plain form flow.
func authErrorReply(c *gin.Context, redirectURL, message string) {
	if wantsSPA(c) {
		c.JSON(http.StatusOK, gin.H{"ok": false, "error": message})
		return
	}
	redirectWithError(c, redirectURL, message)
}

func isSafeRedirectURL(raw string) bool {
	raw = strings.TrimSpace(raw)
	target, err := url.Parse(raw)
	if err != nil || target.User != nil {
		return false
	}
	if target.Scheme == "" {
		return target.Host == "" && !strings.Contains(raw, `\`)
	}
	return (strings.EqualFold(target.Scheme, "http") || strings.EqualFold(target.Scheme, "https")) && target.Host != ""
}
