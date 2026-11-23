package pages

import (
	"context"
	"io"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

// VerifyOTP renders the same auth document as Login but starts on the code
// verification step. It is used for direct visits and as the no-JavaScript
// fallback after the OTP has been sent.
func VerifyOTP(data pages.AuthData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		data.Step = "verify"
		return Layout(data, authCard(data)).Render(ctx, w)
	})
}
