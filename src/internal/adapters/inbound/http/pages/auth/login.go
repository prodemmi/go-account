package pages

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

// Login renders the branded entry page. Both the login and the verification
// steps live in the document; the page starts on the login step.
func Login(data pages.AuthData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		data.Step = "login"
		return Layout(data, authCard(data)).Render(ctx, w)
	})
}
