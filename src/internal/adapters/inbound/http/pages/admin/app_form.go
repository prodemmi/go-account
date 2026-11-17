package admin

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// AppForm renders the admin app create/edit form.
// Generated from admin/app_form.templ by `templ generate`.
func AppForm(app AppFormData, isEdit bool) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		title := "New App"
		if isEdit {
			title = "Edit App"
		}
		action := "/admin/apps"
		if isEdit {
			action = "/admin/apps/" + app.ID + "/update"
		}
		w.Write([]byte(`<div class="header">
	<h1>`))
		pages.WriteEscaped(w, title)
		w.Write([]byte(`</h1>
	<a href="/admin/apps" class="btn-secondary">← Back</a>
</div>
<div class="form-card">
	<form method="POST" action="`))
		pages.WriteEscaped(w, action)
		w.Write([]byte(`">
		<input type="hidden" name="id" value="`))
		pages.WriteEscaped(w, app.ID)
		w.Write([]byte(`" />
		<div class="form-group">
			<label for="name">App Name</label>
			<input type="text" id="name" name="name" value="`))
		pages.WriteEscaped(w, app.Name)
		w.Write([]byte(`" placeholder="Fotomi" required />
		</div>
		<div class="form-group">
			<label for="url">App URL</label>
			<input type="text" id="url" name="url" value="`))
		pages.WriteEscaped(w, app.Url)
		w.Write([]byte(`" placeholder="https://fotomi.net/" required />
		</div>
		<div class="form-group">
			<label for="auth_url">Auth Callback URL</label>
			<input type="text" id="auth_url" name="auth_url" value="`))
		pages.WriteEscaped(w, app.AuthUrl)
		w.Write([]byte(`" placeholder="https://fotomi.net/auth/callback" required />
		</div>
		<div class="form-group">
			<label for="otp_type">OTP Type</label>
			<input type="text" id="otp_type" name="otp_type" value="`))
		pages.WriteEscaped(w, app.OtpType)
		w.Write([]byte(`" placeholder="email" />
		</div>
		<div class="form-group">
			<label for="session_timeout">Session Timeout (seconds)</label>
			<input type="number" id="session_timeout" name="session_timeout" value="`))
		pages.WriteEscaped(w, itostr(app.SessionTimeout))
		w.Write([]byte(`" placeholder="3600" />
		</div>
		<hr style="margin:24px 0; border:none; border-top:1px solid #e2e8f0;" />
		<h3 style="font-size:1rem; margin-bottom:12px; color:#475569;">Auth Page Customization</h3>
		<div class="checkbox-group">
			<input type="checkbox" id="auth_page_enabled" name="auth_page_enabled" value="true"`))
		if app.AuthPageEnabled {
			w.Write([]byte(` checked`))
		}
		w.Write([]byte(` />
			<label for="auth_page_enabled">Enable custom auth page</label>
		</div>
		<div class="form-group">
			<label for="auth_page_url">Custom Auth Page URL</label>
			<input type="text" id="auth_page_url" name="auth_page_url" value="`))
		pages.WriteEscaped(w, app.AuthPageURL)
		w.Write([]byte(`" placeholder="https://fotomi.net/auth/login" />
		</div>
		<div class="form-group">
			<label for="auth_page_logo">Logo URL</label>
			<input type="text" id="auth_page_logo" name="auth_page_logo" value="`))
		pages.WriteEscaped(w, app.AuthPageLogo)
		w.Write([]byte(`" placeholder="https://fotomi.net/logo.png" />
		</div>
		<div class="form-group">
			<label for="auth_page_primary_color">Primary Color</label>
			<input type="color" id="auth_page_primary_color" name="auth_page_primary_color" value="`))
		pages.WriteEscaped(w, app.AuthPagePrimaryColor)
		w.Write([]byte(`" />
		</div>
		<div class="form-group">
			<label for="auth_page_template">Template</label>
			<input type="text" id="auth_page_template" name="auth_page_template" value="`))
		pages.WriteEscaped(w, app.AuthPageTemplate)
		w.Write([]byte(`" placeholder="default" />
		</div>
		<button type="submit" class="btn">Save App</button>
	</form>
</div>`))
		return nil
	})
}

// AppFormData is used by the admin app form template.
type AppFormData struct {
	ID                   string
	Name                 string
	Url                  string
	AuthUrl              string
	OtpType              string
	SessionTimeout       int
	AuthPageURL          string
	AuthPageTemplate     string
	AuthPageLogo         string
	AuthPagePrimaryColor string
	AuthPageEnabled      bool
}
