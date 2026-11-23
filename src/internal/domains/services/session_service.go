package services

import (
	"context"
	"time"

	"aman/internal/domains/event"
	"aman/internal/domains/session"
	"aman/internal/ports/inbound"
	"aman/internal/ports/outbound"
)

// SessionService lets a user list and terminate their own CAS sessions.
// CAS does not manage users: ownership is proven by the app + identifier pair
// carried by the authenticated request.
type SessionService struct {
	sessions outbound.SessionRepository
	apps     outbound.AppRepository
	webhooks *WebhookService
}

func NewSessionService(
	sessions outbound.SessionRepository,
	apps outbound.AppRepository,
	webhooks *WebhookService,
) *SessionService {
	return &SessionService{sessions: sessions, apps: apps, webhooks: webhooks}
}

// List returns the active sessions for an app + identifier pair.
func (s *SessionService) List(_ context.Context, appID, identifier string) ([]session.Session, error) {
	return s.sessions.FindActiveByAppAndIdentifier(appID, identifier)
}

// ListAll returns a filtered, paginated page of sessions (admin panel).
func (s *SessionService) ListAll(_ context.Context, filter inbound.SessionFilter) ([]session.Session, int64, error) {
	return s.sessions.FindPage(outbound.SessionQuery{
		AppID:  filter.AppID,
		Search: filter.Search,
		Status: filter.Status,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

// Stats returns dashboard aggregates for all apps or one app. A nil since
// means "all time"; otherwise only sessions created inside the range count.
func (s *SessionService) Stats(_ context.Context, appID string, since *time.Time) (inbound.SessionStats, error) {
	stats, err := s.sessions.Stats(appID, time.Now(), since)
	if err != nil {
		return inbound.SessionStats{}, err
	}
	return inbound.SessionStats{
		Total:    stats.Total,
		Active:   stats.Active,
		Pending:  stats.Pending,
		Expired:  stats.Expired,
		Revoked:  stats.Revoked,
		Infinite: stats.Infinite,
		Today:    stats.Today,
		Unique:   stats.Unique,
	}, nil
}

// DailyCounts returns the sessions-per-day series used by the dashboard chart.
func (s *SessionService) DailyCounts(_ context.Context, appID string, since time.Time) ([]inbound.DailyCount, error) {
	rows, err := s.sessions.DailyCounts(appID, since)
	if err != nil {
		return nil, err
	}
	counts := make([]inbound.DailyCount, len(rows))
	for i, row := range rows {
		counts[i] = inbound.DailyCount{Day: row.Day, Count: row.Count}
	}
	return counts, nil
}

// Terminate revokes a single session owned by the identifier.
func (s *SessionService) Terminate(ctx context.Context, appID, identifier, id string) error {
	if err := s.ensureOwned(appID, identifier, id); err != nil {
		return err
	}
	if err := s.sessions.Revoke(id); err != nil {
		return err
	}
	s.emitRevoked(ctx, appID, identifier, []string{id})
	return nil
}

// TerminateMany revokes several sessions in one call.
func (s *SessionService) TerminateMany(ctx context.Context, appID, identifier string, ids []string) error {
	for _, id := range ids {
		if err := s.ensureOwned(appID, identifier, id); err != nil {
			return err
		}
	}
	if err := s.sessions.RevokeMany(ids); err != nil {
		return err
	}
	s.emitRevoked(ctx, appID, identifier, ids)
	return nil
}

// TerminateAll revokes every session of an app + identifier pair.
func (s *SessionService) TerminateAll(ctx context.Context, appID, identifier string) error {
	if err := s.sessions.RevokeAll(appID, identifier); err != nil {
		return err
	}
	s.emitRevoked(ctx, appID, identifier, nil)
	return nil
}

// Revoke unconditionally revokes a session (admin panel).
func (s *SessionService) Revoke(_ context.Context, id string) error {
	return s.sessions.Revoke(id)
}

func (s *SessionService) ensureOwned(appID, identifier, id string) error {
	sess, err := s.sessions.FindByID(id)
	if err != nil {
		return err
	}
	if sess.AppID.String() != appID || sess.Identifier != identifier {
		return session.ErrNotFound
	}
	return nil
}

func (s *SessionService) emitRevoked(ctx context.Context, appID, identifier string, ids []string) {
	a, err := s.apps.FindByID(appID)
	if err != nil {
		return
	}
	appEntity := a.Normalize()
	s.webhooks.Emit(ctx, &appEntity, event.SessionRevoked, map[string]any{
		"identifier":  identifier,
		"session_ids": ids,
		"all":         ids == nil,
	})
}
