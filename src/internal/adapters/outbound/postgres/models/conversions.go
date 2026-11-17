package models

import (
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/otp"
	"aman/internal/domains/session"
)

// ToModelApp maps a domain App to its gorm persistence model.
func ToModelApp(a app.App) App {
	return App{
		BaseModel: BaseModel{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			DeletedAt: a.DeletedAt,
		},
		Name:               a.Name,
		Url:                a.Url,
		AuthUrl:            a.AuthUrl,
		OtpType:            a.OtpType,
		SessionTimeout:     int(a.SessionTimeout.Seconds()),
		AuthPageURL:        a.AuthPageURL,
		AuthPageTemplate:   a.AuthPageTemplate,
		AuthPageLogo:       a.AuthPageLogo,
		AuthPagePrimaryColor: a.AuthPagePrimaryColor,
		AuthPageEnabled:    a.AuthPageEnabled,
	}
}

// ToDomainApp maps a gorm App model to its domain entity.
func ToDomainApp(m App) app.App {
	return app.App{
		ID:                 m.ID,
		Name:               m.Name,
		Url:                m.Url,
		AuthUrl:            m.AuthUrl,
		OtpType:            m.OtpType,
		SessionTimeout:     time.Duration(m.SessionTimeout) * time.Second,
		AuthPageURL:        m.AuthPageURL,
		AuthPageTemplate:   m.AuthPageTemplate,
		AuthPageLogo:       m.AuthPageLogo,
		AuthPagePrimaryColor: m.AuthPagePrimaryColor,
		AuthPageEnabled:    m.AuthPageEnabled,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
		DeletedAt:          m.DeletedAt,
	}
}

// ToModelSession maps a domain Session to its gorm persistence model.
func ToModelSession(s session.Session) Session {
	return Session{
		BaseModel: BaseModel{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			DeletedAt: s.DeletedAt,
		},
		AppID:     s.AppID,
		Token:     s.Token,
		ExpiresAt: s.ExpiresAt,
	}
}

// ToDomainSession maps a gorm Session model to its domain entity.
func ToDomainSession(m Session) session.Session {
	return session.Session{
		ID:        m.ID,
		AppID:     m.AppID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// ToModelOtp maps a domain Otp to its gorm persistence model.
func ToModelOtp(o otp.Otp) Otp {
	return Otp{
		BaseModel: BaseModel{
			ID:        o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			DeletedAt: o.DeletedAt,
		},
		AppID:     o.AppID,
		Code:      o.Code,
		Type:      o.Type,
		Used:      o.Used,
		Hash:      o.Hash,
		ExpireAt:  o.ExpireAt,
	}
}

// ToDomainOtp maps a gorm Otp model to its domain entity.
func ToDomainOtp(m Otp) otp.Otp {
	return otp.Otp{
		ID:        m.ID,
		AppID:     m.AppID,
		Type:      m.Type,
		Code:      m.Code,
		Used:      m.Used,
		Hash:      m.Hash,
		ExpireAt:  m.ExpireAt,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}
