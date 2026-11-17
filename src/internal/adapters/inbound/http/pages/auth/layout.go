package auth

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"
	"github.com/a-h/templ"
)

// Layout renders the base HTML layout for auth pages.
// Generated from auth/layout.templ by `templ generate`.
func Layout(title string, content templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
	<title>`))
		pages.WriteEscaped(w, title)
		w.Write([]byte(` - CAS</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet">
	<link rel="stylesheet" href="/static/css/auth.css" />
</head>
<body>`))
		if err := content.Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`
</body>
</html>`))
		return nil
	})
}

