package admin

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// AdminLayout renders the base layout for the admin panel.
// Generated from admin/layout.templ by `templ generate`.
func AdminLayout(title string, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	<title>`))
		pages.WriteEscaped(w, title)
		w.Write([]byte(` - CAS Admin</title>
	<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet">
	<link rel="stylesheet" href="/static/css/admin.css" />
</head>
<body>
	<aside class="sidebar">
		<h2>CAS Admin</h2>
		<nav>
			<a href="/admin/apps" class="active">Apps</a>
			<a href="/admin/sessions">Sessions</a>
		</nav>
	</aside>
	<main class="main">
`))
		if err := content.Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`
	</main>
</body>
</html>`))
		return nil
	})
}
