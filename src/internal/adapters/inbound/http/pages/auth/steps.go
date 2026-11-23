package pages

import (
	"context"
	"io"
	"strconv"
	"strings"

	"aman/internal/adapters/inbound/http/pages"
	"aman/internal/domains/app"

	"github.com/a-h/templ"
)

// authCard renders both steps of the auth flow in one document. CSS toggles
// which step is visible, and the page JavaScript switches between them after
// the OTP is sent, so login and verification no longer need a full reload.
// Without JavaScript each step still has a normal form action and redirects.
func authCard(data pages.AuthData) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		loginHidden := ""
		verifyHidden := " hidden"
		if data.Step == "verify" {
			loginHidden = " hidden"
			verifyHidden = ""
		}

		loginTitle := data.Title
		if loginTitle == "" && !data.CopyCustom {
			loginTitle = "ورود به " + data.App.Name
		}
		loginSubtitle := data.Subtitle
		if loginSubtitle == "" && !data.CopyCustom {
			switch {
			case data.AllowEmail && data.AllowMobile:
				loginSubtitle = "برای دریافت کد ورود، ایمیل یا شماره موبایل خود را وارد کنید."
			case data.AllowEmail:
				loginSubtitle = "برای دریافت کد ورود، ایمیل خود را وارد کنید."
			default:
				loginSubtitle = "برای دریافت کد ورود، شماره موبایل خود را وارد کنید."
			}
		}
		identifierLabel := data.IdentifierLabel
		if identifierLabel == "" && !data.CopyCustom {
			switch {
			case data.AllowEmail && data.AllowMobile:
				identifierLabel = "ایمیل یا شماره موبایل"
			case data.AllowEmail:
				identifierLabel = "ایمیل"
			default:
				identifierLabel = "شماره موبایل"
			}
		}
		sendButton := data.SendButton
		if sendButton == "" {
			sendButton = "دریافت کد ورود"
		}

		verifyTitle := data.VerifyTitle
		if verifyTitle == "" && !data.CopyCustom {
			verifyTitle = "کد تایید را وارد کنید"
		}
		verifySubtitleTemplate := data.VerifySubtitle
		if verifySubtitleTemplate == "" && !data.CopyCustom {
			verifySubtitleTemplate = "کد ارسال‌شده به {identifier} را وارد کنید."
		}
		verifySubtitle := strings.ReplaceAll(verifySubtitleTemplate, "{identifier}", data.Identifier)
		verifyButton := data.VerifyButton
		if verifyButton == "" {
			verifyButton = "تایید و ورود"
		}
		retryText := data.RetryText
		if retryText == "" {
			retryText = "ارسال دوباره کد"
		}
		backText := data.BackText
		if backText == "" {
			backText = "تغییر ایمیل یا شماره موبایل"
		}

		length := data.OTPLength
		if length < 4 {
			length = 5
		}
		if length > 32 {
			length = 32
		}

		logoPosition := data.LogoPosition
		if logoPosition == "" {
			logoPosition = "top"
		}
		if logoPosition == "top" {
			if err := writeBrand(w, data); err != nil {
				return err
			}
		}

		w.Write([]byte(`<div class="card" data-auth-card>`))
		if logoPosition != "top" && logoPosition != "bottom" {
			w.Write([]byte(`<div class="auth-brand-slot">`))
			if err := writeBrand(w, data); err != nil {
				return err
			}
			w.Write([]byte(`</div>`))
		}

		w.Write([]byte(`

	<section class="auth-step" data-auth-step="login"`))
		w.Write([]byte(loginHidden))
		w.Write([]byte(`>`))
		if loginTitle != "" {
			w.Write([]byte(`<h1>`))
			pages.WriteEscaped(w, loginTitle)
			w.Write([]byte(`</h1>`))
		}
		if loginSubtitle != "" {
			w.Write([]byte(`<p class="subtitle">`))
			pages.WriteEscaped(w, loginSubtitle)
			w.Write([]byte(`</p>`))
		}

		w.Write([]byte(`

		<form method="POST" action="/auth/send-otp" data-auth-form data-auth-login>
			<input type="hidden" name="app_id" value="`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`" />`))
		if data.ReturnTo != "" {
			w.Write([]byte(`<input type="hidden" name="redirect" value="`))
			pages.WriteEscaped(w, data.ReturnTo)
			w.Write([]byte(`" />`))
		}
		w.Write([]byte(`<div class="form-group">`))
		if identifierLabel != "" {
			w.Write([]byte(`<label for="identifier">`))
			pages.WriteEscaped(w, identifierLabel)
			w.Write([]byte(`</label>`))
		}
		w.Write([]byte(`<input type="`))
		pages.WriteEscaped(w, data.OTPInputType)
		w.Write([]byte(`" id="identifier" name="identifier" required autofocus autocomplete="username" `))
		switch data.OTPInputType {
		case "email":
			w.Write([]byte(`placeholder="you@example.com"`))
		case "tel":
			w.Write([]byte(`placeholder="09123456789" inputmode="tel"`))
		default:
			w.Write([]byte(`placeholder="ایمیل یا شماره موبایل"`))
		}
		w.Write([]byte(` />
			</div>`))

		if data.CaptchaEnabled {
			w.Write([]byte(`<div class="form-group captcha-group">
				<div class="arcaptcha" data-site-key="`))
			pages.WriteEscaped(w, data.CaptchaSiteKey)
			w.Write([]byte(`" data-lang="fa"`))
			if data.Theme == app.ThemeDark {
				w.Write([]byte(` data-theme="dark"`))
			}
			w.Write([]byte(`></div>
			</div>`))
		}

		w.Write([]byte(`<button type="submit">`))
		pages.WriteEscaped(w, sendButton)
		w.Write([]byte(`</button>
		</form>`))

		if data.FooterText != "" {
			w.Write([]byte(`<p class="page-note">`))
			pages.WriteEscaped(w, data.FooterText)
			w.Write([]byte(`</p>`))
		}
		w.Write([]byte(`</section>

	<section class="auth-step" data-auth-step="verify" data-expires="`))
		if data.OTPExpiresAt > 0 {
			pages.WriteEscaped(w, strconv.FormatInt(data.OTPExpiresAt, 10))
		}
		w.Write([]byte(`"`))
		w.Write([]byte(verifyHidden))
		w.Write([]byte(`>`))
		if verifyTitle != "" {
			w.Write([]byte(`<h1>`))
			pages.WriteEscaped(w, verifyTitle)
			w.Write([]byte(`</h1>`))
		}
		if verifySubtitleTemplate != "" {
			w.Write([]byte(`<p class="subtitle" data-verify-subtitle data-template="`))
			pages.WriteEscaped(w, verifySubtitleTemplate)
			w.Write([]byte(`">`))
			pages.WriteEscaped(w, verifySubtitle)
			w.Write([]byte(`</p>`))
		}

		w.Write([]byte(`

		<form method="POST" action="/auth/verify-otp" data-auth-form data-auth-verify>
			<input type="hidden" name="app_id" value="`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`" />
			<input type="hidden" name="identifier" value="`))
		pages.WriteEscaped(w, data.Identifier)
		w.Write([]byte(`" />
			<input type="hidden" name="otp_hash" value="`))
		pages.WriteEscaped(w, data.OtpHash)
		w.Write([]byte(`" />`))
		if data.OTPExpiresAt > 0 {
			w.Write([]byte(`<input type="hidden" name="expires_at" value="`))
			pages.WriteEscaped(w, strconv.FormatInt(data.OTPExpiresAt, 10))
			w.Write([]byte(`" />`))
		} else {
			w.Write([]byte(`<input type="hidden" name="expires_at" value="" />`))
		}
		if data.ReturnTo != "" {
			w.Write([]byte(`<input type="hidden" name="redirect" value="`))
			pages.WriteEscaped(w, data.ReturnTo)
			w.Write([]byte(`" />`))
		}
		w.Write([]byte(`<div class="pin" data-pin style="--pin-count:`))
		pages.WriteEscaped(w, strconv.Itoa(length))
		w.Write([]byte(`">
			<input type="hidden" name="otp_code" data-pin-value />`))

		for i := 0; i < length; i++ {
			w.Write([]byte(`<input class="pin-input" type="`))
			pages.WriteEscaped(w, data.OTPCodeInputType)
			w.Write([]byte(`" maxlength="1" autocomplete="`))
			if i == 0 {
				w.Write([]byte(`one-time-code`))
			} else {
				w.Write([]byte(`off`))
			}
			w.Write([]byte(`" aria-label="رقم `))
			pages.WriteEscaped(w, strconv.Itoa(i+1))
			w.Write([]byte(`"`))
			if i == 0 {
				w.Write([]byte(` autofocus`))
			}
			if data.OTPCodeInputType == "tel" {
				w.Write([]byte(` inputmode="numeric" pattern="[0-9]*"`))
			} else {
				w.Write([]byte(` autocapitalize="off" spellcheck="false"`))
			}
			w.Write([]byte(` />`))
		}

		w.Write([]byte(`</div>
			<button type="submit">`))
		pages.WriteEscaped(w, verifyButton)
		w.Write([]byte(`</button>
		</form>

		<form method="POST" action="/auth/send-otp" class="retry-form" data-auth-retry>
			<input type="hidden" name="app_id" value="`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`" />
			<input type="hidden" name="identifier" value="`))
		pages.WriteEscaped(w, data.Identifier)
		w.Write([]byte(`" />
			<input type="hidden" name="retry" value="1" />`))
		if data.ReturnTo != "" {
			w.Write([]byte(`<input type="hidden" name="redirect" value="`))
			pages.WriteEscaped(w, data.ReturnTo)
			w.Write([]byte(`" />`))
		}
		w.Write([]byte(`<button type="submit" class="btn-secondary" data-retry-button>`))
		pages.WriteEscaped(w, retryText)
		w.Write([]byte(`<span class="retry-countdown" data-retry-countdown></span></button>
		</form>

		<div class="footer">
			<a href="/auth/login?app_id=`))
		pages.WriteEscaped(w, data.AppID)
		if data.ReturnTo != "" {
			w.Write([]byte(`&redirect=`))
			pages.WriteEscaped(w, data.ReturnTo)
		}
		w.Write([]byte(`" data-auth-back>`))
		pages.WriteEscaped(w, backText)
		w.Write([]byte(`</a>
		</div>`))

		if data.FooterText != "" {
			w.Write([]byte(`<p class="page-note">`))
			pages.WriteEscaped(w, data.FooterText)
			w.Write([]byte(`</p>`))
		}

		w.Write([]byte(`</section>`))
		if logoPosition == "bottom" {
			if err := writeBrand(w, data); err != nil {
				return err
			}
		}
		w.Write([]byte(`</div>`))
		return nil
	})
}

// writeBrand renders the logo and app name using the configured mask and
// position. It is shared by both auth steps so the header never jumps when the
// page switches from login to verification.
func writeBrand(w io.Writer, data pages.AuthData) error {
	position := data.LogoPosition
	if position == "" {
		position = "top"
	}
	mask := data.LogoMask
	if mask == "" {
		mask = "none"
	}

	if position == "top" {
		w.Write([]byte(`<div class="auth-brand-outside">`))
	}
	w.Write([]byte(`<div class="brand brand-`))
	pages.WriteEscaped(w, position)
	w.Write([]byte(`">`))
	if data.LogoURL != "" {
		w.Write([]byte(`<span class="logo logo-mask-`))
		pages.WriteEscaped(w, mask)
		w.Write([]byte(`" data-skeleton><span class="media-skeleton" aria-hidden="true"></span><img src="`))
		pages.WriteEscaped(w, data.LogoURL)
		w.Write([]byte(`" alt="`))
		pages.WriteEscaped(w, data.App.Name)
		w.Write([]byte(`" decoding="async" /></span>`))
	} else {
		w.Write([]byte(`<span class="logo brand-initial logo-mask-`))
		pages.WriteEscaped(w, mask)
		w.Write([]byte(`">`))
		pages.WriteEscaped(w, firstLetter(data.App.Name))
		w.Write([]byte(`</span>`))
	}
	w.Write([]byte(`<div class="name">`))
	pages.WriteEscaped(w, data.App.Name)
	w.Write([]byte(`</div></div>`))
	if position == "top" {
		w.Write([]byte(`</div>`))
	}
	return nil
}

func firstLetter(name string) string {
	for _, r := range name {
		return string(r)
	}
	return "A"
}
