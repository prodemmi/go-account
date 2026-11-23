package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"aman/internal/domains/app"
)

var errInvalidIdentifier = errors.New("identifier is required")

// resolveChannel validates the identifier against the app's OTP channel and
// returns the channel that must be used to deliver the code.
func resolveChannel(a app.App, identifier string) (string, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return "", errInvalidIdentifier
	}

	isEmail := strings.Contains(identifier, "@")

	switch a.OTPChannel {
	case "", app.OTPChannelSMS:
		if isEmail {
			return "", fmt.Errorf("%w: this app only accepts a mobile number", app.ErrInvalid)
		}
		return app.OTPChannelSMS, nil
	case app.OTPChannelEmail:
		if !isEmail {
			return "", fmt.Errorf("%w: this app only accepts an email address", app.ErrInvalid)
		}
		return app.OTPChannelEmail, nil
	case app.OTPChannelBoth:
		if isEmail {
			return app.OTPChannelEmail, nil
		}
		return app.OTPChannelSMS, nil
	default:
		return "", app.ErrInvalidChannel
	}
}

func splitIdentifier(identifier string) (mobile, email string) {
	if strings.Contains(identifier, "@") {
		return "", identifier
	}
	return identifier, ""
}

// appendQuery adds a query parameter to a URL without losing existing ones.
func appendQuery(rawURL, key, value string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		sep := "?"
		if strings.Contains(rawURL, "?") {
			sep = "&"
		}
		return rawURL + sep + key + "=" + url.QueryEscape(value)
	}

	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u.String()
}

// randomToken returns a URL-safe opaque token.
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
