package services

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/event"
	"aman/internal/domains/otp"
	"aman/internal/domains/ott"
	"aman/internal/domains/session"
	"aman/internal/ports/inbound"
	"aman/internal/ports/outbound"
)

// AuthService owns the CAS authentication flow: send OTP, verify OTP, check
// an existing CAS session, exchange an OTT for credentials and log out.
type AuthService struct {
	apps      outbound.AppRepository
	sessions  outbound.SessionRepository
	otps      outbound.OTPRepository
	otts      outbound.OTTRepository
	hasher    outbound.AuthHash
	jwt       outbound.JWTService
	generator outbound.OTPGenerator
	notifier  outbound.Notifier
	webhooks  *WebhookService
}

func NewAuthService(
	apps outbound.AppRepository,
	sessions outbound.SessionRepository,
	otps outbound.OTPRepository,
	otts outbound.OTTRepository,
	hasher outbound.AuthHash,
	jwt outbound.JWTService,
	generator outbound.OTPGenerator,
	notifier outbound.Notifier,
	webhooks *WebhookService,
) *AuthService {
	return &AuthService{
		apps:      apps,
		sessions:  sessions,
		otps:      otps,
		otts:      otts,
		hasher:    hasher,
		jwt:       jwt,
		generator: generator,
		notifier:  notifier,
		webhooks:  webhooks,
	}
}

// SendOTP creates an OTP for the given identifier and delivers it through the
// channel configured on the app (sms, email or both).
func (s *AuthService) SendOTP(ctx context.Context, cmd inbound.SendOTPCommand) (*inbound.OTPChallenge, error) {
	a, err := s.apps.FindByID(cmd.AppID)
	if err != nil {
		return nil, err
	}
	appEntity := a.Normalize()

	identifier := strings.TrimSpace(cmd.Identifier)
	channel, err := resolveChannel(appEntity, identifier)
	if err != nil {
		return nil, err
	}

	code, err := s.generator.GenerateOTP(appEntity.OTPLength, appEntity.OTPCharset)
	if err != nil {
		return nil, err
	}

	hash := s.hasher.CreateAuthHash(appEntity, cmd.UserAgent)
	expiresAt := time.Now().Add(appEntity.OTPExpire)

	if _, err := s.otps.Create(otp.Otp{
		AppID:      appEntity.ID,
		Identifier: identifier,
		Type:       channel,
		Code:       code,
		Hash:       &hash,
		ExpireAt:   &expiresAt,
	}); err != nil {
		return nil, err
	}

	if err := s.deliver(channel, identifier, code); err != nil {
		return nil, err
	}

	s.webhooks.Emit(ctx, &appEntity, event.OTPSent, map[string]any{
		"identifier": identifier,
		"channel":    channel,
		"expires_at": expiresAt,
	})

	return &inbound.OTPChallenge{
		AppID:      appEntity.ID.String(),
		Identifier: identifier,
		Channel:    channel,
		OTPHash:    hash,
		Length:     appEntity.OTPLength,
		ExpiresAt:  expiresAt,
	}, nil
}

// VerifyOTP consumes a valid OTP, creates a pending CAS session and issues an
// OTT the app can exchange for credentials. The session becomes active only
// after VerifyOTT succeeds.
func (s *AuthService) VerifyOTP(ctx context.Context, cmd inbound.VerifyOTPCommand) (*inbound.AuthResult, error) {
	a, err := s.apps.FindByID(cmd.AppID)
	if err != nil {
		return nil, err
	}
	appEntity := a.Normalize()

	if !s.hasher.VerifyAuthHash(cmd.OTPHash, appEntity, cmd.UserAgent) {
		return nil, fmt.Errorf("%w: cannot verify otp hash", otp.ErrInvalid)
	}

	consumed, err := s.otps.Consume(appEntity.ID.String(), strings.TrimSpace(cmd.Identifier), cmd.Code, cmd.OTPHash)
	if err != nil {
		return nil, err
	}

	identifier := consumed.Identifier
	mobile, email := splitIdentifier(identifier)
	now := time.Now()
	ttl := appEntity.SessionTTL()

	token, err := s.jwt.CreateToken(appEntity.ID.String(), identifier, ttl)
	if err != nil {
		return nil, err
	}

	// LastLoginAt is the creation time of the most recent session that already
	// existed before this one. It is not "now": the current login is the row
	// being inserted and its CreatedAt already carries that information.
	var lastLoginAt *time.Time
	if latest, err := s.sessions.LatestCreatedAt(appEntity.ID.String(), identifier); err == nil {
		lastLoginAt = latest
	}

	created, err := s.sessions.Create(session.Session{
		AppID:       appEntity.ID,
		Token:       token,
		Identifier:  identifier,
		Mobile:      mobile,
		Email:       email,
		UserAgent:   cmd.UserAgent,
		IP:          cmd.IP,
		Country:     cmd.Country,
		City:        cmd.City,
		Region:      cmd.Region,
		Device:      cmd.Device,
		Status:      session.StatusPending,
		IsLoggedIn:  false,
		Infinite:    appEntity.InfiniteSession,
		LastLoginAt: lastLoginAt,
		ExpiresAt:   now.Add(ttl),
	})
	if err != nil {
		return nil, err
	}

	ottToken, err := s.issueOTT(&appEntity, &created)
	if err != nil {
		return nil, err
	}

	s.webhooks.Emit(ctx, &appEntity, event.SessionCreated, map[string]any{
		"session_id": created.ID.String(),
		"identifier": identifier,
		"status":     created.Status,
		"expires_at": created.ExpiresAt,
		"infinite":   created.Infinite,
	})

	return &inbound.AuthResult{
		App:         &appEntity,
		Session:     &created,
		OTT:         ottToken,
		RedirectURL: redirectWithOTT(&appEntity, cmd.ReturnTo, ottToken),
	}, nil
}

// Check reports whether the browser still has a valid CAS session. When it
// does, a fresh OTT is issued for the callback redirect.
func (s *AuthService) Check(ctx context.Context, cmd inbound.CheckCommand) (*inbound.CheckResult, error) {
	if strings.TrimSpace(cmd.SessionToken) == "" {
		return &inbound.CheckResult{}, nil
	}

	sess, err := s.sessions.FindByToken(cmd.SessionToken)
	if err != nil || sess.AppID.String() != cmd.AppID || !sess.Active(time.Now()) || !sess.MatchesAgent(cmd.UserAgent) {
		return &inbound.CheckResult{}, nil
	}

	a, err := s.apps.FindByID(cmd.AppID)
	if err != nil {
		return nil, err
	}
	appEntity := a.Normalize()

	ottToken, err := s.issueOTT(&appEntity, sess)
	if err != nil {
		return nil, err
	}

	s.webhooks.Emit(ctx, &appEntity, event.OTTIssued, map[string]any{
		"session_id": sess.ID.String(),
		"identifier": sess.Identifier,
	})

	return &inbound.CheckResult{
		Authenticated: true,
		OTT:           ottToken,
		RedirectURL:   redirectWithOTT(&appEntity, cmd.ReturnTo, ottToken),
	}, nil
}

// VerifyOTT exchanges a one-time ticket for the session credentials and marks
// the pending session as active.
func (s *AuthService) VerifyOTT(ctx context.Context, cmd inbound.VerifyOTTCommand) (*inbound.Credentials, error) {
	ticket, err := s.otts.FindByToken(strings.TrimSpace(cmd.OTT))
	if err != nil {
		return nil, err
	}
	if ticket.AppID.String() != cmd.AppID {
		return nil, ott.ErrForbidden
	}
	if ticket.Used {
		return nil, ott.ErrUsed
	}
	if !ticket.ExpireAt.After(time.Now()) {
		return nil, ott.ErrExpired
	}

	sess, err := s.sessions.FindByID(ticket.SessionID.String())
	if err != nil {
		return nil, err
	}
	if sess.AppID.String() != cmd.AppID {
		return nil, session.ErrForbidden
	}
	if !sess.MatchesAgent(cmd.UserAgent) {
		return nil, session.ErrForbidden
	}
	now := time.Now()
	if !sess.Pending() && !sess.Active(now) {
		return nil, session.ErrExpired
	}

	if err := s.otts.MarkUsed(ticket.ID.String()); err != nil {
		return nil, err
	}
	if err := s.sessions.MarkActive(sess.ID.String()); err != nil {
		return nil, err
	}
	sess.Status = session.StatusActive
	sess.IsLoggedIn = true

	a, err := s.apps.FindByID(cmd.AppID)
	if err != nil {
		return nil, err
	}
	appEntity := a.Normalize()

	s.webhooks.Emit(ctx, &appEntity, event.OTTVerified, map[string]any{
		"session_id": sess.ID.String(),
		"identifier": sess.Identifier,
	})

	return &inbound.Credentials{
		AppID:        appEntity.ID.String(),
		SessionID:    sess.ID.String(),
		SessionToken: sess.Token,
		Identifier:   sess.Identifier,
		Mobile:       sess.Mobile,
		Email:        sess.Email,
		ExpiresAt:    sess.ExpiresAt,
		Infinite:     sess.Infinite,
	}, nil
}

// Logout revokes the current CAS session.
func (s *AuthService) Logout(ctx context.Context, cmd inbound.LogoutCommand) error {
	if strings.TrimSpace(cmd.SessionToken) == "" {
		return nil
	}

	sess, err := s.sessions.FindByToken(cmd.SessionToken)
	if err != nil {
		return nil
	}
	if err := s.sessions.Revoke(sess.ID.String()); err != nil {
		return err
	}

	if a, err := s.apps.FindByID(sess.AppID.String()); err == nil {
		appEntity := a.Normalize()
		s.webhooks.Emit(ctx, &appEntity, event.LoggedOut, map[string]any{
			"session_id": sess.ID.String(),
			"identifier": sess.Identifier,
		})
	}
	return nil
}

func (s *AuthService) issueOTT(a *app.App, sess *session.Session) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}

	if _, err := s.otts.Create(ott.OTT{
		AppID:      a.ID,
		SessionID:  sess.ID,
		Token:      token,
		Identifier: sess.Identifier,
		ExpireAt:   time.Now().Add(a.OTTExpire),
	}); err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) deliver(channel, identifier, code string) error {
	switch channel {
	case app.OTPChannelEmail:
		return s.notifier.SendEmail(identifier, "Your verification code", "Your verification code is "+code)
	default:
		return s.notifier.SendSMS(identifier, "Your verification code is "+code)
	}
}

// redirectWithOTT appends the ticket to the app callback, or to the URL the
// user originally visited when the app allows it and the origin matches.
func redirectWithOTT(a *app.App, returnTo, ottToken string) string {
	target := a.CallbackURL
	if a.RedirectBack && sameOrigin(a.Url, returnTo) {
		target = returnTo
	}
	return appendQuery(target, "ott", ottToken)
}

func sameOrigin(appURL, target string) bool {
	appURL = strings.TrimSpace(appURL)
	target = strings.TrimSpace(target)
	if appURL == "" || target == "" {
		return false
	}

	a, err := url.Parse(appURL)
	if err != nil || a.Scheme == "" || a.Host == "" {
		return false
	}
	t, err := url.Parse(target)
	if err != nil || t.Scheme == "" || t.Host == "" {
		return false
	}
	return strings.EqualFold(a.Scheme, t.Scheme) && strings.EqualFold(a.Host, t.Host)
}
