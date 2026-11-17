package auth

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// Login renders the branded login page for an app.
// Generated from auth/login.templ by `templ generate`.
func Login(data pages.AuthData) templ.Component {
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
		<h1>Sign in to `))
		pages.WriteEscaped(w, data.App.Name)
		w.Write([]byte(`</h1>`))
		if data.Error != "" {
			w.Write([]byte(`<div class="error show">`))
			pages.WriteEscaped(w, data.Error)
			w.Write([]byte(`</div>`))
		}
		if data.Success != "" {
			w.Write([]byte(`<div class="success show">`))
			pages.WriteEscaped(w, data.Success)
			w.Write([]byte(`</div>`))
		}
		w.Write([]byte(`
		<form method="POST" action="/auth/login">
			<input type="hidden" name="app_id" value="`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`" />
			<input type="hidden" name="csrf_token" value="`))
		pages.WriteEscaped(w, data.CsrfToken)
		w.Write([]byte(`" />
			<div class="form-group">
				<label for="email">Email</label>
				<input type="email" id="email" name="email" placeholder="you@example.com" required />
			</div>
			<div class="form-group">
				<label for="password">Password</label>
				<input type="password" id="password" name="password" placeholder="••••••••" required />
			</div>
			<button type="submit">Sign in</button>
		</form>
		<div class="footer">
			New to `))
		pages.WriteEscaped(w, data.App.Name)
		w.Write([]byte(`? <a href="/auth/signup?app_id=`))
		pages.WriteEscaped(w, data.AppID)
		w.Write([]byte(`">Create an account</a>
		</div>
	</div>
</div>`))
		return nil
	})
}
