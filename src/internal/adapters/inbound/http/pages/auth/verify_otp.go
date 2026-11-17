package auth

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// VerifyOTP renders the branded verify OTP page for an app.
// Generated from auth/verify_otp.templ by `templ generate`.
func VerifyOTP(data pages.AuthData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<div class="container">
	<div class="card">
		<div class="brand">`))
		if data.AuthPageLogo != "" {
			w.Write([]byte(`<img src="`))
			pages.WriteEscaped(w, data.AuthPageLogo)
			w.Write([]byte(`" alt="Logo" />`))
		}
		w.Write([]byte(`
			<div class="name" style="color: `))
		pages.WriteEscaped(w, data.AuthPagePrimaryColor)
		w.Write([]byte(`">`))
		pages.WriteEscaped(w, data.App.Name)
		w.Write([]byte(`</div>
		</div>
		<h1>Verify your identity</h1>
		<p style="text-align:center; color:#6b7280; margin-bottom:24px;">A code was sent to your `))
		pages.WriteEscaped(w, data.App.OtpType)
		w.Write([]byte(`.</p>`))
		if data.Error != "" {
			w.Write([]byte(`<div class="error show">`))
			pages.WriteEscaped(w, data.Error)
			w.Write([]byte(`</div>`))
		}
		w.Write([]byte(`
		<form method="POST" action="/auth/verify-otp">
			<input type="hidden" name="app_id" value="`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`" />
			<input type="hidden" name="user_id" value="`))
		pages.WriteEscaped(w, data.App.ID.String())
		w.Write([]byte(`" />
			<input type="hidden" name="otp_type" value="`))
		pages.WriteEscaped(w, data.App.OtpType)
		w.Write([]byte(`" />
			<input type="hidden" name="otp_hash" value="`))
		pages.WriteEscaped(w, data.OtpHash)
		w.Write([]byte(`" />
			<input type="hidden" name="csrf_token" value="`))
		pages.WriteEscaped(w, data.CsrfToken)
		w.Write([]byte(`" />
			<div class="form-group">
				<label for="otp_code">Verification Code</label>
				<input type="text" id="otp_code" name="otp_code" placeholder="12345" maxlength="5" required />
			</div>
			<button type="submit">Verify</button>
		</form>
		<div class="footer">
			Didn't receive a code? <a href="/auth/login?app_id=`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`">Try signing in again</a>
		</div>
	</div>
</div>`))
		return nil
	})
}
