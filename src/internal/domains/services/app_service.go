package services

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"aman/internal/domains/app"
	"aman/internal/ports/inbound"
	"aman/internal/ports/outbound"

	"github.com/google/uuid"
)

// AppService exposes app configuration to the admin panel.
type AppService struct {
	apps outbound.AppRepository
}

func NewAppService(apps outbound.AppRepository) *AppService {
	return &AppService{apps: apps}
}

func (s *AppService) List(context.Context) ([]app.App, error) {
	apps, err := s.apps.FindAll()
	if err != nil {
		return nil, err
	}
	for i := range apps {
		apps[i] = apps[i].Normalize()
	}
	return apps, nil
}

func (s *AppService) ListPage(_ context.Context, filter inbound.AppFilter) ([]app.App, int64, error) {
	apps, total, err := s.apps.FindPage(strings.TrimSpace(filter.Search), filter.CreatedFrom, filter.CreatedUntil, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	for i := range apps {
		apps[i] = apps[i].Normalize()
	}
	return apps, total, nil
}

func (s *AppService) Get(_ context.Context, id string) (*app.App, error) {
	a, err := s.apps.FindByID(id)
	if err != nil {
		return nil, err
	}
	normalized := a.Normalize()
	return &normalized, nil
}

func (s *AppService) Create(_ context.Context, in inbound.AppInput) (app.App, error) {
	a, err := appFromInput(uuid.New(), in)
	if err != nil {
		return app.App{}, err
	}
	return s.apps.Create(a)
}

func (s *AppService) Update(_ context.Context, id string, in inbound.AppInput) (app.App, error) {
	existing, err := s.apps.FindByID(id)
	if err != nil {
		return app.App{}, err
	}

	updated, err := appFromInput(existing.ID, in)
	if err != nil {
		return app.App{}, err
	}
	updated.CreatedAt = existing.CreatedAt
	return s.apps.Update(updated)
}

func (s *AppService) Delete(_ context.Context, id string) error {
	return s.apps.Delete(id)
}

func appFromInput(id uuid.UUID, in inbound.AppInput) (app.App, error) {
	if strings.TrimSpace(in.Name) == "" {
		return app.App{}, errors.New("نام برنامه الزامی است")
	}
	if strings.TrimSpace(in.CallbackURL) == "" {
		return app.App{}, errors.New("آدرس بازگشت (callback) الزامی است")
	}
	if in.OTPLength < 4 || in.OTPLength > 32 {
		return app.App{}, errors.New("طول کد یکبار مصرف باید بین ۴ تا ۳۲ باشد")
	}
	allowedCountries, err := parseAllowedCountries(in.AllowedCountries)
	if err != nil {
		return app.App{}, err
	}
	captchaProvider := strings.ToLower(strings.TrimSpace(in.AuthPageCaptchaProvider))
	if captchaProvider == "" {
		captchaProvider = app.CaptchaProviderARCaptcha
	}
	if captchaProvider != app.CaptchaProviderARCaptcha && captchaProvider != app.CaptchaProviderGoogle {
		return app.App{}, errors.New("ارائه‌دهنده کپچا معتبر نیست")
	}
	textColor, err := authPageColor(in.AuthPageTextColor)
	if err != nil {
		return app.App{}, errors.New("رنگ متن صفحه ورود معتبر نیست")
	}
	mutedTextColor, err := authPageColor(in.AuthPageMutedTextColor)
	if err != nil {
		return app.App{}, errors.New("رنگ متن ثانویه صفحه ورود معتبر نیست")
	}

	a := app.App{
		ID:          id,
		Name:        in.Name,
		Url:         in.URL,
		CallbackURL: in.CallbackURL,
		WebhookURL:  in.WebhookURL,

		OTPChannel: in.OTPChannel,
		OTPLength:  in.OTPLength,
		OTPCharset: in.OTPCharset,
		OTPExpire:  time.Duration(in.OTPExpireSeconds) * time.Second,
		OTTExpire:  time.Duration(in.OTTExpireSeconds) * time.Second,

		SessionTimeout:        time.Duration(in.SessionTimeoutSeconds) * time.Second,
		InfiniteSession:       in.InfiniteSession,
		SessionCookieName:     in.SessionCookieName,
		SessionCookieHTTPOnly: in.SessionCookieHTTPOnly,
		SessionCookieSecure:   in.SessionCookieSecure,
		SessionCookieSameSite: in.SessionCookieSameSite,
		RedirectBack:          in.RedirectBack,

		CustomAuthEnabled:           in.CustomAuthEnabled,
		AllowedCountries:            allowedCountries,
		AuthPageEnabled:             in.AuthPageEnabled,
		AuthPageURL:                 in.AuthPageURL,
		AuthPageTitle:               in.AuthPageTitle,
		AuthPageLayout:              in.AuthPageLayout,
		AuthPageFont:                in.AuthPageFont,
		AuthPageTextAlign:           in.AuthPageTextAlign,
		AuthPageLogo:                in.AuthPageLogo,
		AuthPageLogoMask:            in.AuthPageLogoMask,
		AuthPageLogoPosition:        in.AuthPageLogoPosition,
		AuthPagePrimaryColor:        in.AuthPagePrimaryColor,
		AuthPageBackgroundColor:     in.AuthPageBackgroundColor,
		AuthPageCardBackgroundColor: in.AuthPageCardBackgroundColor,
		AuthPageTextColor:           textColor,
		AuthPageMutedTextColor:      mutedTextColor,
		AuthPagePinBackgroundColor:  in.AuthPagePinBackgroundColor,
		AuthPagePinTextColor:        in.AuthPagePinTextColor,
		AuthPageTheme:               in.AuthPageTheme,
		AuthPageCaptchaEnabled:      in.AuthPageCaptchaEnabled,
		AuthPageCaptchaProvider:     captchaProvider,
		AuthPageCaptchaSiteKey:      strings.TrimSpace(in.AuthPageCaptchaSiteKey),
		AuthPageCaptchaSecretKey:    strings.TrimSpace(in.AuthPageCaptchaSecretKey),
		AuthPageBannerMode:          in.AuthPageBannerMode,
		AuthPageBanners:             in.AuthPageBanners,
		AuthPageCopyCustom:          in.AuthPageCopyCustom,
		AuthPageLoginTitle:          in.AuthPageLoginTitle,
		AuthPageLoginSubtitle:       in.AuthPageLoginSubtitle,
		AuthPageIdentifierLabel:     in.AuthPageIdentifierLabel,
		AuthPageSendButton:          in.AuthPageSendButton,
		AuthPageVerifyTitle:         in.AuthPageVerifyTitle,
		AuthPageVerifySubtitle:      in.AuthPageVerifySubtitle,
		AuthPageVerifyButton:        in.AuthPageVerifyButton,
		AuthPageRetryText:           in.AuthPageRetryText,
		AuthPageBackText:            in.AuthPageBackText,
		AuthPageFooterText:          in.AuthPageFooterText,
	}
	return a.Normalize(), nil
}

func authPageColor(raw string) (string, error) {
	color := strings.TrimSpace(raw)
	if color == "" {
		return "", nil
	}
	if len(color) != 7 || color[0] != '#' {
		return "", errors.New("invalid hex color")
	}
	if _, err := hex.DecodeString(color[1:]); err != nil {
		return "", err
	}
	return color, nil
}

func parseAllowedCountries(raw string) ([]string, error) {
	seen := make(map[string]struct{})
	countries := make([]string, 0)
	for _, value := range strings.Split(raw, ",") {
		country := strings.ToUpper(strings.TrimSpace(value))
		if country == "" {
			continue
		}
		if len(country) != 2 || country[0] < 'A' || country[0] > 'Z' || country[1] < 'A' || country[1] > 'Z' {
			return nil, errors.New("کشورها را با کد دوحرفی ISO، جداشده با ویرگول وارد کنید")
		}
		if _, exists := seen[country]; exists {
			continue
		}
		seen[country] = struct{}{}
		countries = append(countries, country)
	}
	return countries, nil
}
