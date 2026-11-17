package admin

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// SessionsList renders the admin sessions list page.
// Generated from admin/sessions.templ by `templ generate`.
func SessionsList(sessions []SessionRow) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<div class="header">
	<h1>Active Sessions</h1>
</div>`))
		if len(sessions) == 0 {
			w.Write([]byte(`<div class="empty-state">No active sessions.</div>`))
		} else {
			w.Write([]byte(`<table class="table">
				<thead>
					<tr>
						<th>App ID</th>
						<th>Token</th>
						<th>Expires At</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>`))
			for _, s := range sessions {
				w.Write([]byte(`<tr>
						<td>`))
				pages.WriteEscaped(w, s.AppID)
				w.Write([]byte(`</td>
						<td style="font-family:monospace; font-size:0.875rem;">`))
				pages.WriteEscaped(w, s.Token)
				w.Write([]byte(`</td>
						<td>`))
				pages.WriteEscaped(w, s.ExpiresAt)
				w.Write([]byte(`</td>
						<td>
							<form method="POST" action="/admin/sessions/`))
				pages.WriteEscaped(w, s.Token)
				w.Write([]byte(`/delete" style="display:inline;" onsubmit="return confirm('Revoke this session?');">
								<button type="submit" class="btn-secondary" style="background:#ef4444;">Revoke</button>
							</form>
						</td>
					</tr>`))
			}
			w.Write([]byte(`</tbody></table>`))
		}
		return nil
	})
}

// SessionRow is used by the admin sessions list template.
type SessionRow struct {
	AppID     string
	Token     string
	ExpiresAt string
}
