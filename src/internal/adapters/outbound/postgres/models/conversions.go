package models

import (
	"encoding/json"
	"strings"
	"time"

	"aman/internal/domains/app"
	"aman/internal/domains/otp"
	"aman/internal/domains/ott"
	"aman/internal/domains/session"
)

func encodeBanners(banners []app.Banner) string {
	if len(banners) == 0 {
		return "[]"
	}
	raw, err := json.Marshal(banners)
	if err != nil {
		return "[]"
	}
	return string(raw)
}

func decodeBanners(raw string) []app.Banner {
	if raw == "" {
		return nil
	}
	var banners []app.Banner
	if err := json.Unmarshal([]byte(raw), &banners); err != nil {
		return nil
	}
	return banners
}

// ToModelApp maps a domain App to its gorm persistence model.
func ToModelApp(a app.App) App {
	return App{
		BaseModel: BaseModel{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			DeletedAt: a.DeletedAt,
		},
		Name:                        a.Name,
		Url:                         a.Url,
		CallbackURL:                 a.CallbackURL,
		WebhookURL:                  a.WebhookURL,
		OTPChannel:                  a.OTPChannel,
		OTPLength:                   a.OTPLength,
		OTPCharset:                  a.OTPCharset,
		OTPExpire:                   int(a.OTPExpire.Seconds()),
		OTTExpire:                   int(a.OTTExpire.Seconds()),
		SessionTimeout:              int(a.SessionTimeout.Seconds()),
		InfiniteSession:             a.InfiniteSession,
		SessionCookieName:           a.SessionCookieName,
		SessionCookieHTTPOnly:       a.SessionCookieHTTPOnly,
		SessionCookieSecure:         a.SessionCookieSecure,
		SessionCookieSameSite:       a.SessionCookieSameSite,
		RedirectBack:                a.RedirectBack,
		CustomAuthEnabled:           a.CustomAuthEnabled,
		AllowedCountries:            strings.Join(a.AllowedCountries, ","),
		AuthPageEnabled:             a.AuthPageEnabled,
		AuthPageURL:                 a.AuthPageURL,
		AuthPageTitle:               a.AuthPageTitle,
		AuthPageLayout:              a.AuthPageLayout,
		AuthPageFont:                a.AuthPageFont,
		AuthPageTextAlign:           a.AuthPageTextAlign,
		AuthPageLogo:                a.AuthPageLogo,
		AuthPageLogoMask:            a.AuthPageLogoMask,
		AuthPageLogoPosition:        a.AuthPageLogoPosition,
		AuthPagePrimaryColor:        a.AuthPagePrimaryColor,
		AuthPageBackgroundColor:     a.AuthPageBackgroundColor,
		AuthPageCardBackgroundColor: a.AuthPageCardBackgroundColor,
		AuthPageTextColor:           a.AuthPageTextColor,
		AuthPageMutedTextColor:      a.AuthPageMutedTextColor,
		AuthPagePinBackgroundColor:  a.AuthPagePinBackgroundColor,
		AuthPagePinTextColor:        a.AuthPagePinTextColor,
		AuthPageTheme:               a.AuthPageTheme,
		AuthPageCaptchaEnabled:      a.AuthPageCaptchaEnabled,
		AuthPageCaptchaProvider:     a.AuthPageCaptchaProvider,
		AuthPageCaptchaSiteKey:      a.AuthPageCaptchaSiteKey,
		AuthPageCaptchaSecretKey:    a.AuthPageCaptchaSecretKey,
		AuthPageBannerMode:          a.AuthPageBannerMode,
		AuthPageBanners:             encodeBanners(a.AuthPageBanners),
		AuthPageCopyCustom:          a.AuthPageCopyCustom,
		AuthPageLoginTitle:          a.AuthPageLoginTitle,
		AuthPageLoginSubtitle:       a.AuthPageLoginSubtitle,
		AuthPageIdentifierLabel:     a.AuthPageIdentifierLabel,
		AuthPageSendButton:          a.AuthPageSendButton,
		AuthPageVerifyTitle:         a.AuthPageVerifyTitle,
		AuthPageVerifySubtitle:      a.AuthPageVerifySubtitle,
		AuthPageVerifyButton:        a.AuthPageVerifyButton,
		AuthPageRetryText:           a.AuthPageRetryText,
		AuthPageBackText:            a.AuthPageBackText,
		AuthPageFooterText:          a.AuthPageFooterText,
	}
}

// ToDomainApp maps a gorm App model to its domain entity.
func ToDomainApp(m App) app.App {
	return app.App{
		ID:                          m.ID,
		Name:                        m.Name,
		Url:                         m.Url,
		CallbackURL:                 m.CallbackURL,
		WebhookURL:                  m.WebhookURL,
		OTPChannel:                  m.OTPChannel,
		OTPLength:                   m.OTPLength,
		OTPCharset:                  m.OTPCharset,
		OTPExpire:                   time.Duration(m.OTPExpire) * time.Second,
		OTTExpire:                   time.Duration(m.OTTExpire) * time.Second,
		SessionTimeout:              time.Duration(m.SessionTimeout) * time.Second,
		InfiniteSession:             m.InfiniteSession,
		SessionCookieName:           m.SessionCookieName,
		SessionCookieHTTPOnly:       m.SessionCookieHTTPOnly,
		SessionCookieSecure:         m.SessionCookieSecure,
		SessionCookieSameSite:       m.SessionCookieSameSite,
		RedirectBack:                m.RedirectBack,
		CustomAuthEnabled:           m.CustomAuthEnabled,
		AllowedCountries:            splitCountries(m.AllowedCountries),
		AuthPageEnabled:             m.AuthPageEnabled,
		AuthPageURL:                 m.AuthPageURL,
		AuthPageTitle:               m.AuthPageTitle,
		AuthPageLayout:              m.AuthPageLayout,
		AuthPageFont:                m.AuthPageFont,
		AuthPageTextAlign:           m.AuthPageTextAlign,
		AuthPageLogo:                m.AuthPageLogo,
		AuthPageLogoMask:            m.AuthPageLogoMask,
		AuthPageLogoPosition:        m.AuthPageLogoPosition,
		AuthPagePrimaryColor:        m.AuthPagePrimaryColor,
		AuthPageBackgroundColor:     m.AuthPageBackgroundColor,
		AuthPageCardBackgroundColor: m.AuthPageCardBackgroundColor,
		AuthPageTextColor:           m.AuthPageTextColor,
		AuthPageMutedTextColor:      m.AuthPageMutedTextColor,
		AuthPagePinBackgroundColor:  m.AuthPagePinBackgroundColor,
		AuthPagePinTextColor:        m.AuthPagePinTextColor,
		AuthPageTheme:               m.AuthPageTheme,
		AuthPageCaptchaEnabled:      m.AuthPageCaptchaEnabled,
		AuthPageCaptchaProvider:     m.AuthPageCaptchaProvider,
		AuthPageCaptchaSiteKey:      m.AuthPageCaptchaSiteKey,
		AuthPageCaptchaSecretKey:    m.AuthPageCaptchaSecretKey,
		AuthPageBannerMode:          m.AuthPageBannerMode,
		AuthPageBanners:             decodeBanners(m.AuthPageBanners),
		AuthPageCopyCustom:          m.AuthPageCopyCustom,
		AuthPageLoginTitle:          m.AuthPageLoginTitle,
		AuthPageLoginSubtitle:       m.AuthPageLoginSubtitle,
		AuthPageIdentifierLabel:     m.AuthPageIdentifierLabel,
		AuthPageSendButton:          m.AuthPageSendButton,
		AuthPageVerifyTitle:         m.AuthPageVerifyTitle,
		AuthPageVerifySubtitle:      m.AuthPageVerifySubtitle,
		AuthPageVerifyButton:        m.AuthPageVerifyButton,
		AuthPageRetryText:           m.AuthPageRetryText,
		AuthPageBackText:            m.AuthPageBackText,
		AuthPageFooterText:          m.AuthPageFooterText,
		CreatedAt:                   m.CreatedAt,
		UpdatedAt:                   m.UpdatedAt,
		DeletedAt:                   m.DeletedAt,
	}
}

func splitCountries(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	return strings.Split(raw, ",")
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
		AppID:       s.AppID,
		Token:       s.Token,
		Identifier:  s.Identifier,
		Mobile:      s.Mobile,
		Email:       s.Email,
		UserAgent:   s.UserAgent,
		IP:          s.IP,
		Country:     s.Country,
		City:        s.City,
		Region:      s.Region,
		Device:      s.Device,
		Status:      s.Status,
		IsLoggedIn:  s.IsLoggedIn,
		Infinite:    s.Infinite,
		LastLoginAt: s.LastLoginAt,
		ExpiresAt:   s.ExpiresAt,
	}
}

// ToDomainSession maps a gorm Session model to its domain entity.
func ToDomainSession(m Session) session.Session {
	return session.Session{
		ID:          m.ID,
		AppID:       m.AppID,
		Token:       m.Token,
		Identifier:  m.Identifier,
		Mobile:      m.Mobile,
		Email:       m.Email,
		UserAgent:   m.UserAgent,
		IP:          m.IP,
		Country:     m.Country,
		City:        m.City,
		Region:      m.Region,
		Device:      m.Device,
		Status:      m.Status,
		IsLoggedIn:  m.IsLoggedIn,
		Infinite:    m.Infinite,
		LastLoginAt: m.LastLoginAt,
		ExpiresAt:   m.ExpiresAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
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
		AppID:      o.AppID,
		Identifier: o.Identifier,
		Code:       o.Code,
		Type:       o.Type,
		Used:       o.Used,
		Hash:       o.Hash,
		ExpireAt:   o.ExpireAt,
	}
}

// ToDomainOtp maps a gorm Otp model to its domain entity.
func ToDomainOtp(m Otp) otp.Otp {
	return otp.Otp{
		ID:         m.ID,
		AppID:      m.AppID,
		Identifier: m.Identifier,
		Type:       m.Type,
		Code:       m.Code,
		Used:       m.Used,
		Hash:       m.Hash,
		ExpireAt:   m.ExpireAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
	}
}

// ToModelOTT maps a domain OTT to its gorm persistence model.
func ToModelOTT(o ott.OTT) OTT {
	return OTT{
		BaseModel: BaseModel{
			ID:        o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: o.UpdatedAt,
			DeletedAt: o.DeletedAt,
		},
		AppID:      o.AppID,
		SessionID:  o.SessionID,
		Token:      o.Token,
		Identifier: o.Identifier,
		Used:       o.Used,
		ExpireAt:   o.ExpireAt,
	}
}

// ToDomainOTT maps a gorm OTT model to its domain entity.
func ToDomainOTT(m OTT) ott.OTT {
	return ott.OTT{
		ID:         m.ID,
		AppID:      m.AppID,
		SessionID:  m.SessionID,
		Token:      m.Token,
		Identifier: m.Identifier,
		Used:       m.Used,
		ExpireAt:   m.ExpireAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
	}
}
