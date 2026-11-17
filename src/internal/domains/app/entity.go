package app

import (
	"time"

	"github.com/google/uuid"
)

// App is the domain entity representing a client application in the CAS.
// Each app owns its own session timeout and can customize its auth pages.
type App struct {
	ID             uuid.UUID
	Name           string
	Url            string
	AuthUrl        string
	OtpType        string
	SessionTimeout time.Duration // session time per app

	// Auth page customization (CAS serves branded auth pages per app)
	AuthPageURL         string // custom auth page URL (if app hosts its own)
	AuthPageTemplate    string // template name to render (default "default")
	AuthPageLogo        string // logo URL for branding
	AuthPagePrimaryColor string // primary brand color (hex)
	AuthPageEnabled     bool   // whether custom auth page is enabled

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
