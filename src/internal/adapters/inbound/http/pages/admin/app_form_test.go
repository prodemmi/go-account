package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"aman/internal/domains/app"
)

func TestAppFormContextualGuideLinks(t *testing.T) {
	var form bytes.Buffer
	if err := AppForm(AppFormData{AuthPageBanners: []app.Banner{{Duration: 6}}}, false).Render(context.Background(), &form); err != nil {
		t.Fatal(err)
	}
	html := form.String()

	for _, field := range []string{
		"name", "url", "callback_url", "webhook_url",
		"otp_channel", "otp_length", "otp_charset", "otp_expire", "ott_expire",
		"session_timeout", "session_cookie_name", "session_cookie_same_site", "infinite_session",
		"session_cookie_http_only", "session_cookie_secure", "redirect_back", "allowed_countries",
		"custom_auth_enabled", "auth_page_enabled", "auth_page_url", "auth_page_layout", "auth_page_font",
		"auth_page_text_align", "auth_page_theme", "auth_page_primary_color", "auth_page_background_color",
		"auth_page_card_background_color", "auth_page_text_color", "auth_page_muted_text_color",
		"auth_page_pin_background_color", "auth_page_pin_text_color", "auth_page_logo", "auth_page_logo_mask",
		"auth_page_logo_position", "auth_page_captcha_enabled", "auth_page_captcha_provider",
		"auth_page_captcha_site_key", "auth_page_captcha_secret_key", "auth_page_banner_mode",
		"banner_url", "banner_text", "banner_link", "banner_duration", "auth_page_title",
		"auth_page_login_title", "auth_page_login_subtitle", "auth_page_identifier_label", "auth_page_send_button",
		"auth_page_verify_title", "auth_page_verify_subtitle", "auth_page_verify_button", "auth_page_retry_text",
		"auth_page_back_text", "auth_page_footer_text",
	} {
		if !strings.Contains(html, `name="`+field+`"`) {
			t.Errorf("settings form is missing field %q", field)
		}
	}

	var guide bytes.Buffer
	if err := Guide().Render(context.Background(), &guide); err != nil {
		t.Fatal(err)
	}
	for _, anchor := range []string{
		"app-basics", "webhook", "otp", "sessions", "geo-restrictions", "auth-page", "branding", "captcha",
	} {
		if !strings.Contains(html, `data-guide-anchor="`+anchor+`"`) {
			t.Errorf("settings form has no guide mapping for #%s", anchor)
		}
		if !strings.Contains(guide.String(), `id="`+anchor+`"`) {
			t.Errorf("guide has no destination for #%s", anchor)
		}
	}

	for _, want := range []string{
		`document.querySelectorAll('#app-form .field, #app-form .switch-row')`,
		`field.closest('[data-guide-anchor]')`,
		`link.href = '/admin/guide#'`,
		`link.target = '_blank'`,
		`link.rel = 'noopener noreferrer'`,
		`href="/admin/guide#branding" target="_blank" rel="noopener noreferrer"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("contextual guide links are missing %q", want)
		}
	}
}
