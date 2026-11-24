package postgres

import (
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/otp"
	"aman/internal/domains/ott"
	"aman/internal/domains/session"
	"aman/internal/ports/outbound"
)

// Function-field mocks for the postgres repositories. They live next to the
// adapter so tests can swap the database for an in-memory double.
//
// A nil function returns the zero value and a nil error.

type MockAppRepository struct {
	CreateFn   func(a app.App) (app.App, error)
	FindByIDFn func(id string) (*app.App, error)
	FindAllFn  func() ([]app.App, error)
	FindPageFn func(search string, createdFrom, createdUntil *time.Time, limit, offset int) ([]app.App, int64, error)
	UpdateFn   func(a app.App) (app.App, error)
	DeleteFn   func(id string) error
}

func (m *MockAppRepository) Create(a app.App) (app.App, error) {
	if m.CreateFn != nil {
		return m.CreateFn(a)
	}
	return a, nil
}

func (m *MockAppRepository) FindByID(id string) (*app.App, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, app.ErrNotFound
}

func (m *MockAppRepository) FindAll() ([]app.App, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn()
	}
	return nil, nil
}

func (m *MockAppRepository) FindPage(search string, createdFrom, createdUntil *time.Time, limit, offset int) ([]app.App, int64, error) {
	if m.FindPageFn != nil {
		return m.FindPageFn(search, createdFrom, createdUntil, limit, offset)
	}
	return nil, 0, nil
}

func (m *MockAppRepository) Update(a app.App) (app.App, error) {
	if m.UpdateFn != nil {
		return m.UpdateFn(a)
	}
	return a, nil
}

func (m *MockAppRepository) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

type MockSessionRepository struct {
	CreateFn                       func(s session.Session) (session.Session, error)
	FindByIDFn                     func(id string) (*session.Session, error)
	FindByTokenFn                  func(token string) (*session.Session, error)
	FindAllFn                      func() ([]session.Session, error)
	FindPageFn                     func(q outbound.SessionQuery) ([]session.Session, int64, error)
	FindActiveByAppAndIdentifierFn func(appID, identifier string) ([]session.Session, error)
	LatestCreatedAtFn              func(appID, identifier string) (*time.Time, error)
	StatsFn                        func(appID, status string, now time.Time, since, until *time.Time) (outbound.SessionStats, error)
	DailyCountsFn                  func(appID, status string, since time.Time, until *time.Time) ([]outbound.DailyCount, error)
	MarkActiveFn                   func(id string) error
	RevokeFn                       func(id string) error
	RevokeManyFn                   func(ids []string) error
	RevokeAllFn                    func(appID, identifier string) error
	DeleteExpiredFn                func(before time.Time) error
}

func (m *MockSessionRepository) Create(s session.Session) (session.Session, error) {
	if m.CreateFn != nil {
		return m.CreateFn(s)
	}
	return s, nil
}

func (m *MockSessionRepository) FindByID(id string) (*session.Session, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, session.ErrNotFound
}

func (m *MockSessionRepository) FindByToken(token string) (*session.Session, error) {
	if m.FindByTokenFn != nil {
		return m.FindByTokenFn(token)
	}
	return nil, session.ErrNotFound
}

func (m *MockSessionRepository) FindAll() ([]session.Session, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn()
	}
	return nil, nil
}

func (m *MockSessionRepository) FindPage(q outbound.SessionQuery) ([]session.Session, int64, error) {
	if m.FindPageFn != nil {
		return m.FindPageFn(q)
	}
	return nil, 0, nil
}

func (m *MockSessionRepository) FindActiveByAppAndIdentifier(appID, identifier string) ([]session.Session, error) {
	if m.FindActiveByAppAndIdentifierFn != nil {
		return m.FindActiveByAppAndIdentifierFn(appID, identifier)
	}
	return nil, nil
}

func (m *MockSessionRepository) LatestCreatedAt(appID, identifier string) (*time.Time, error) {
	if m.LatestCreatedAtFn != nil {
		return m.LatestCreatedAtFn(appID, identifier)
	}
	return nil, nil
}

func (m *MockSessionRepository) Stats(appID, status string, now time.Time, since, until *time.Time) (outbound.SessionStats, error) {
	if m.StatsFn != nil {
		return m.StatsFn(appID, status, now, since, until)
	}
	return outbound.SessionStats{}, nil
}

func (m *MockSessionRepository) DailyCounts(appID, status string, since time.Time, until *time.Time) ([]outbound.DailyCount, error) {
	if m.DailyCountsFn != nil {
		return m.DailyCountsFn(appID, status, since, until)
	}
	return nil, nil
}

func (m *MockSessionRepository) MarkActive(id string) error {
	if m.MarkActiveFn != nil {
		return m.MarkActiveFn(id)
	}
	return nil
}

func (m *MockSessionRepository) Revoke(id string) error {
	if m.RevokeFn != nil {
		return m.RevokeFn(id)
	}
	return nil
}

func (m *MockSessionRepository) RevokeMany(ids []string) error {
	if m.RevokeManyFn != nil {
		return m.RevokeManyFn(ids)
	}
	return nil
}

func (m *MockSessionRepository) RevokeAll(appID, identifier string) error {
	if m.RevokeAllFn != nil {
		return m.RevokeAllFn(appID, identifier)
	}
	return nil
}

func (m *MockSessionRepository) DeleteExpired(before time.Time) error {
	if m.DeleteExpiredFn != nil {
		return m.DeleteExpiredFn(before)
	}
	return nil
}

type MockOTPRepository struct {
	CreateFn  func(o otp.Otp) (otp.Otp, error)
	ConsumeFn func(appID, identifier, code, hash string) (otp.Otp, error)
}

func (m *MockOTPRepository) Create(o otp.Otp) (otp.Otp, error) {
	if m.CreateFn != nil {
		return m.CreateFn(o)
	}
	return o, nil
}

func (m *MockOTPRepository) Consume(appID, identifier, code, hash string) (otp.Otp, error) {
	if m.ConsumeFn != nil {
		return m.ConsumeFn(appID, identifier, code, hash)
	}
	return otp.Otp{}, otp.ErrNotFound
}

type MockOTTRepository struct {
	CreateFn        func(o ott.OTT) (ott.OTT, error)
	FindByTokenFn   func(token string) (*ott.OTT, error)
	MarkUsedFn      func(id string) error
	DeleteExpiredFn func(before time.Time) error
}

func (m *MockOTTRepository) Create(o ott.OTT) (ott.OTT, error) {
	if m.CreateFn != nil {
		return m.CreateFn(o)
	}
	return o, nil
}

func (m *MockOTTRepository) FindByToken(token string) (*ott.OTT, error) {
	if m.FindByTokenFn != nil {
		return m.FindByTokenFn(token)
	}
	return nil, ott.ErrNotFound
}

func (m *MockOTTRepository) MarkUsed(id string) error {
	if m.MarkUsedFn != nil {
		return m.MarkUsedFn(id)
	}
	return nil
}

func (m *MockOTTRepository) DeleteExpired(before time.Time) error {
	if m.DeleteExpiredFn != nil {
		return m.DeleteExpiredFn(before)
	}
	return nil
}
