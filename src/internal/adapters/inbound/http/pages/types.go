package pages

import (
	"io"

	"aman/internal/domains/app"
)

// AuthData holds the data passed to auth page templates.
type AuthData struct {
	App              app.App
	AppID            string
	Error            string
	Success          string
	OtpHash          string
	CsrfToken        string
	PageTitle        string
	PrimaryColor     string
	LogoURL          string
	AuthPageLogo     string
	AuthPagePrimaryColor string
}

// WriteEscaped writes HTML-escaped text to w.
func WriteEscaped(w io.Writer, s string) {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			w.Write([]byte("&amp;"))
		case '<':
			w.Write([]byte("&lt;"))
		case '>':
			w.Write([]byte("&gt;"))
		case '"':
			w.Write([]byte("&quot;"))
		case '\'':
			w.Write([]byte("&#39;"))
		default:
			w.Write([]byte{s[i]})
		}
	}
}