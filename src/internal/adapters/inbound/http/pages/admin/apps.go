package admin

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// AppsList renders the admin apps list page.
// Generated from admin/apps.templ by `templ generate`.
func AppsList(apps []AppRow) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<div class="header">
	<h1>Applications</h1>
	<a href="/admin/apps/new" class="btn">+ Add App</a>
</div>`))
		if len(apps) == 0 {
			w.Write([]byte(`<div class="empty-state">No apps registered yet.</div>`))
		} else {
			w.Write([]byte(`<table class="table">
				<thead>
					<tr>
						<th>Name</th>
						<th>URL</th>
						<th>Otp Type</th>
						<th>Session Timeout</th>
						<th>Auth Page</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>`))
			for _, app := range apps {
				w.Write([]byte(`<tr>
						<td><strong>`))
				pages.WriteEscaped(w, app.Name)
				w.Write([]byte(`</strong></td>
						<td>`))
				pages.WriteEscaped(w, app.Url)
				w.Write([]byte(`</td>
						<td>`))
				pages.WriteEscaped(w, app.OtpType)
				w.Write([]byte(`</td>
						<td>`))
				pages.WriteEscaped(w, itostr(app.SessionTimeout))
				w.Write([]byte(`s</td>
						<td>`))
				if app.AuthPageEnabled {
					w.Write([]byte(`<span class="badge badge-success">Custom</span>`))
				} else {
					w.Write([]byte(`<span class="badge badge-default">Default</span>`))
				}
				w.Write([]byte(`</td>
						<td>
							<a href="/admin/apps/`))
				pages.WriteEscaped(w, app.ID)
				w.Write([]byte(`" class="btn-secondary" style="padding:6px 10px;">Edit</a>
							<form method="POST" action="/admin/apps/`))
				pages.WriteEscaped(w, app.ID)
				w.Write([]byte(`/delete" style="display:inline;" onsubmit="return confirm('Delete this app?');">
								<button type="submit" class="btn-secondary" style="background:#ef4444;">Delete</button>
							</form>
						</td>
					</tr>`))
			}
			w.Write([]byte(`</tbody></table>`))
		}
		return nil
	})
}

func itostr(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

// AppRow is used by the admin apps list template.
type AppRow struct {
	ID              string
	Name            string
	Url             string
	OtpType         string
	SessionTimeout  int
	AuthPageEnabled bool
}
