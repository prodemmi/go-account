package admin

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"aman/internal/adapters/inbound/http/pages"
	adminpages "aman/internal/adapters/inbound/http/pages/admin"
	"aman/internal/domains/app"
	"aman/internal/domains/session"
	"aman/internal/ports/inbound"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

const (
	notFoundPath = "/not-found"

	appsPageSize     = 10
	sessionsPageSize = 15
)

// AdminHandler manages apps and sessions through the server-rendered admin
// panel. There is no admin JSON API.
type AdminHandler struct {
	apps     inbound.AppUseCase
	sessions inbound.SessionUseCase
}

func NewAdminHandler(apps inbound.AppUseCase, sessions inbound.SessionUseCase) *AdminHandler {
	return &AdminHandler{apps: apps, sessions: sessions}
}

// AppsList renders the admin apps list with server-side search and paging.
func (h AdminHandler) AppsList(c *gin.Context) {
	page, offset := pageParams(c, appsPageSize)
	search := strings.TrimSpace(c.Query("q"))
	from, until, dateFrom, dateTo, filterError := parseDateRange(c)

	apps, total, err := h.apps.ListPage(c.Request.Context(), inbound.AppFilter{
		Search:       search,
		CreatedFrom:  from,
		CreatedUntil: until,
		Limit:        appsPageSize,
		Offset:       offset,
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	rows := make([]adminpages.AppRow, len(apps))
	for i, a := range apps {
		rows[i] = adminpages.AppRow{
			ID:              a.ID.String(),
			Name:            a.Name,
			Url:             a.Url,
			CallbackURL:     a.CallbackURL,
			OTPChannel:      otpChannelLabel(a.OTPChannel),
			OTPExpire:       int(a.OTPExpire.Seconds()),
			OTTExpire:       int(a.OTTExpire.Seconds()),
			SessionTimeout:  int(a.SessionTimeout.Seconds()),
			InfiniteSession: a.InfiniteSession,
			AuthPageEnabled: a.AuthPageEnabled,
			CreatedAt:       adminpages.FormatDateTime(a.CreatedAt),
		}
	}

	meta := adminpages.PageMeta{
		BasePath:    "/admin/apps",
		Search:      search,
		From:        dateFrom,
		To:          dateTo,
		FilterError: filterError,
		Page:        page,
		Pages:       totalPages(total, appsPageSize),
		Total:       total,
		PageSize:    appsPageSize,
	}

	render(c, "برنامه‌ها", "/admin/apps", adminpages.AppsList(rows, meta, c.Query("error"), c.Query("success")))
}

// NewApp renders the new app form.
func (h AdminHandler) NewApp(c *gin.Context) {
	render(c, "برنامه جدید", "/admin/apps", adminpages.AppForm(adminpages.AppFormData{
		OTPChannel:                  app.OTPChannelSMS,
		OTPLength:                   5,
		OTPCharset:                  app.OTPCharsetNumeric,
		OTPExpire:                   120,
		OTTExpire:                   60,
		SessionTimeout:              3600,
		SessionCookieName:           app.DefaultSessionCookieName,
		SessionCookieHTTPOnly:       true,
		SessionCookieSameSite:       app.SameSiteLax,
		AuthPageTheme:               app.ThemeLight,
		AuthPageLayout:              app.LayoutCard,
		AuthPageFont:                app.FontVazirmatn,
		AuthPageTextAlign:           app.AlignCenter,
		AuthPageLogoMask:            app.LogoMaskNone,
		AuthPageLogoPosition:        app.LogoPositionTop,
		AuthPagePrimaryColor:        "#3b82f6",
		AuthPageBackgroundColor:     "#f8fafc",
		AuthPageCardBackgroundColor: "#ffffff",
		AuthPageTextColor:           "#111827",
		AuthPageMutedTextColor:      "#667085",
		AuthPagePinBackgroundColor:  "#f9fafb",
		AuthPagePinTextColor:        "#111827",
		AuthPageBannerMode:          app.BannerModeNone,
		Error:                       c.Query("error"),
	}, false))
}

// CreateApp creates a new app from the admin form.
func (h AdminHandler) CreateApp(c *gin.Context) {
	created, err := h.apps.Create(c.Request.Context(), parseAppInput(c))
	if err != nil {
		redirectError(c, "/admin/apps/new", err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/admin/apps/"+created.ID.String()+"?success=created")
}

// EditApp renders the edit app form.
func (h AdminHandler) EditApp(c *gin.Context) {
	a, err := h.apps.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		redirectNotFound(c)
		return
	}

	// Show the effective copy in the inputs rather than empty fields. Saving
	// the form marks the copy as custom, so empty would otherwise permanently
	// hide the built-in strings after any unrelated edit.
	copyValue := func(current, fallback string) string {
		if a.AuthPageCopyCustom || strings.TrimSpace(current) != "" {
			return current
		}
		return fallback
	}
	loginSubtitle := "برای دریافت کد ورود، شماره موبایل خود را وارد کنید."
	identifierLabel := "شماره موبایل"
	switch {
	case a.AllowsEmail() && a.AllowsMobile():
		loginSubtitle = "برای دریافت کد ورود، ایمیل یا شماره موبایل خود را وارد کنید."
		identifierLabel = "ایمیل یا شماره موبایل"
	case a.AllowsEmail():
		loginSubtitle = "برای دریافت کد ورود، ایمیل خود را وارد کنید."
		identifierLabel = "ایمیل"
	}

	render(c, "ویرایش برنامه", "/admin/apps", adminpages.AppForm(adminpages.AppFormData{
		ID:                          a.ID.String(),
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
		AuthPageBanners:             a.AuthPageBanners,
		AuthPageCopyCustom:          a.AuthPageCopyCustom,
		AuthPageLoginTitle:          copyValue(a.AuthPageLoginTitle, "ورود به {app}"),
		AuthPageLoginSubtitle:       copyValue(a.AuthPageLoginSubtitle, loginSubtitle),
		AuthPageIdentifierLabel:     copyValue(a.AuthPageIdentifierLabel, identifierLabel),
		AuthPageSendButton:          copyValue(a.AuthPageSendButton, "دریافت کد ورود"),
		AuthPageVerifyTitle:         copyValue(a.AuthPageVerifyTitle, "کد تایید را وارد کنید"),
		AuthPageVerifySubtitle:      copyValue(a.AuthPageVerifySubtitle, "کد ارسال‌شده به {identifier} را وارد کنید."),
		AuthPageVerifyButton:        copyValue(a.AuthPageVerifyButton, "تایید و ورود"),
		AuthPageRetryText:           copyValue(a.AuthPageRetryText, "ارسال دوباره کد"),
		AuthPageBackText:            copyValue(a.AuthPageBackText, pages.EditIdentifierText(a.AllowsEmail(), a.AllowsMobile())),
		AuthPageFooterText:          copyValue(a.AuthPageFooterText, "ورود امن با Aman CAS"),
		Error:                       c.Query("error"),
		Success:                     c.Query("success"),
	}, true))
}

// UpdateApp updates an existing app and returns to the edit form so the admin
// can keep working. The save result is shown as a notification, not a redirect
// back to the apps list.
func (h AdminHandler) UpdateApp(c *gin.Context) {
	id := c.Param("id")
	if _, err := h.apps.Update(c.Request.Context(), id, parseAppInput(c)); err != nil {
		redirectError(c, "/admin/apps/"+id, err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/admin/apps/"+id+"?success=updated")
}

// DeleteApp deletes an app.
func (h AdminHandler) DeleteApp(c *gin.Context) {
	if err := h.apps.Delete(c.Request.Context(), c.Param("id")); err != nil {
		redirectError(c, "/admin/apps", err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/admin/apps?success=deleted")
}

// Guide renders the written admin documentation.
func (h AdminHandler) Guide(c *gin.Context) {
	render(c, "راهنما", "/admin/guide", adminpages.Guide())
}

// Developers renders the API reference for app developers.
func (h AdminHandler) Developers(c *gin.Context) {
	render(c, "توسعه‌دهندگان", "/admin/developers", adminpages.Developers())
}

// Dashboard renders the stats dashboard with an app filter and charts.
func (h AdminHandler) Dashboard(c *gin.Context) {
	appID := strings.TrimSpace(c.Query("app_id"))
	status := strings.TrimSpace(c.Query("status"))
	switch status {
	case session.StatusActive, session.StatusPending, session.StatusRevoked, "expired":
	default:
		status = ""
	}
	days, allTime := parseRange(c.Query("days"))
	from, until, dateFrom, dateTo, filterError := parseDateRange(c)
	customDateRange := filterError == "" && (dateFrom != "" || dateTo != "")

	apps, err := h.apps.List(c.Request.Context())
	if err != nil {
		redirectNotFound(c)
		return
	}

	now := time.Now()
	var since *time.Time
	if customDateRange {
		since = from
	} else if !allTime {
		start := startOfDay(now).AddDate(0, 0, -(days - 1))
		since = &start
	}

	stats, err := h.sessions.Stats(c.Request.Context(), appID, status, since, until)
	if err != nil {
		redirectNotFound(c)
		return
	}

	// An all-time chart starts at the first recorded day. For the regular
	// ranges it starts at the first day of the selected window.
	chartSince := time.Time{}
	if since != nil {
		chartSince = *since
	}
	counts, err := h.sessions.DailyCounts(c.Request.Context(), appID, status, chartSince, until)
	if err != nil {
		redirectNotFound(c)
		return
	}

	chartEnd := startOfDay(now)
	if until != nil {
		chartEnd = until.AddDate(0, 0, -1)
	}
	chartDays := days
	if customDateRange {
		chartStart := chartStart(chartSince, counts, chartEnd)
		chartDays = int(chartEnd.Sub(chartStart).Hours()/24) + 1
		if chartDays < 1 {
			chartDays = 1
		}
	} else if allTime {
		chartDays = daysBetween(counts, chartEnd)
	}

	render(c, "داشبورد", "/admin", adminpages.Dashboard(adminpages.DashboardData{
		Apps:        toAppOptions(apps),
		AppID:       appID,
		Status:      status,
		Days:        days,
		AllTime:     allTime,
		From:        dateFrom,
		To:          dateTo,
		FilterError: filterError,
		Stats:       stats,
		Daily:       fillDaily(counts, chartStart(chartSince, counts, chartEnd), chartDays),
		ChartDays:   chartDays,
		Error:       c.Query("error"),
	}))
}

// SessionsList renders the admin sessions list with server-side filters.
func (h AdminHandler) SessionsList(c *gin.Context) {
	page, offset := pageParams(c, sessionsPageSize)
	search := strings.TrimSpace(c.Query("q"))
	status := strings.TrimSpace(c.Query("status"))
	appID := strings.TrimSpace(c.Query("app_id"))
	from, until, dateFrom, dateTo, filterError := parseDateRange(c)

	sessions, total, err := h.sessions.ListAll(c.Request.Context(), inbound.SessionFilter{
		AppID:        appID,
		Search:       search,
		Status:       status,
		CreatedFrom:  from,
		CreatedUntil: until,
		Limit:        sessionsPageSize,
		Offset:       offset,
	})
	if err != nil {
		redirectNotFound(c)
		return
	}

	apps, err := h.apps.List(c.Request.Context())
	if err != nil {
		redirectNotFound(c)
		return
	}
	appNames := make(map[string]string, len(apps))
	for _, a := range apps {
		appNames[a.ID.String()] = a.Name
	}

	rows := make([]adminpages.SessionRow, len(sessions))
	now := time.Now()
	for i, s := range sessions {
		lastLogin := "—"
		if s.LastLoginAt != nil {
			lastLogin = adminpages.FormatDateTime(*s.LastLoginAt)
		}
		appName := appNames[s.AppID.String()]
		if appName == "" {
			appName = "برنامه حذف‌شده"
		}
		rows[i] = adminpages.SessionRow{
			ID:          s.ID.String(),
			AppID:       s.AppID.String(),
			AppName:     appName,
			Identifier:  s.Identifier,
			IP:          s.IP,
			Country:     s.Country,
			City:        s.City,
			Region:      s.Region,
			Device:      s.Device,
			UserAgent:   s.UserAgent,
			LastLoginAt: lastLogin,
			CreatedAt:   adminpages.FormatDateTime(s.CreatedAt),
			ExpiresAt:   adminpages.FormatDateTime(s.ExpiresAt),
			Status:      s.DisplayStatus(now),
			Infinite:    s.Infinite,
		}
	}

	meta := adminpages.PageMeta{
		BasePath:    "/admin/sessions",
		Search:      search,
		Status:      status,
		AppID:       appID,
		From:        dateFrom,
		To:          dateTo,
		FilterError: filterError,
		Page:        page,
		Pages:       totalPages(total, sessionsPageSize),
		Total:       total,
		PageSize:    sessionsPageSize,
	}

	render(c, "نشست‌ها", "/admin/sessions", adminpages.SessionsList(rows, meta, toAppOptions(apps), c.Query("error"), c.Query("success")))
}

// DeleteSession revokes a session.
func (h AdminHandler) DeleteSession(c *gin.Context) {
	if err := h.sessions.Revoke(c.Request.Context(), c.Param("id")); err != nil {
		redirectError(c, "/admin/sessions", err.Error())
		return
	}
	c.Redirect(http.StatusFound, "/admin/sessions?success=revoked")
}

func parseAppInput(c *gin.Context) inbound.AppInput {
	otpLength := atoiDefault(c.PostForm("otp_length"), 5)
	if otpLength < 4 {
		otpLength = 4
	}
	if otpLength > 32 {
		otpLength = 32
	}

	cookieName := strings.TrimSpace(c.PostForm("session_cookie_name"))
	if cookieName == "" {
		cookieName = app.DefaultSessionCookieName
	}

	return inbound.AppInput{
		Name:        c.PostForm("name"),
		URL:         c.PostForm("url"),
		CallbackURL: c.PostForm("callback_url"),
		WebhookURL:  c.PostForm("webhook_url"),

		OTPChannel:       c.PostForm("otp_channel"),
		OTPLength:        otpLength,
		OTPCharset:       c.PostForm("otp_charset"),
		OTPExpireSeconds: atoiDefault(c.PostForm("otp_expire"), 120),
		OTTExpireSeconds: atoiDefault(c.PostForm("ott_expire"), 60),

		SessionTimeoutSeconds: atoiDefault(c.PostForm("session_timeout"), 3600),
		InfiniteSession:       checkbox(c.PostForm("infinite_session")),
		SessionCookieName:     cookieName,
		SessionCookieHTTPOnly: checkbox(c.PostForm("session_cookie_http_only")),
		SessionCookieSecure:   checkbox(c.PostForm("session_cookie_secure")),
		SessionCookieSameSite: c.PostForm("session_cookie_same_site"),
		RedirectBack:          checkbox(c.PostForm("redirect_back")),

		CustomAuthEnabled:           checkbox(c.PostForm("custom_auth_enabled")),
		AllowedCountries:            c.PostForm("allowed_countries"),
		AuthPageEnabled:             checkbox(c.PostForm("auth_page_enabled")),
		AuthPageURL:                 c.PostForm("auth_page_url"),
		AuthPageTitle:               c.PostForm("auth_page_title"),
		AuthPageLayout:              c.PostForm("auth_page_layout"),
		AuthPageFont:                c.PostForm("auth_page_font"),
		AuthPageTextAlign:           c.PostForm("auth_page_text_align"),
		AuthPageLogo:                c.PostForm("auth_page_logo"),
		AuthPageLogoMask:            c.PostForm("auth_page_logo_mask"),
		AuthPageLogoPosition:        c.PostForm("auth_page_logo_position"),
		AuthPagePrimaryColor:        c.PostForm("auth_page_primary_color"),
		AuthPageBackgroundColor:     c.PostForm("auth_page_background_color"),
		AuthPageCardBackgroundColor: c.PostForm("auth_page_card_background_color"),
		AuthPageTextColor:           c.PostForm("auth_page_text_color"),
		AuthPageMutedTextColor:      c.PostForm("auth_page_muted_text_color"),
		AuthPagePinBackgroundColor:  c.PostForm("auth_page_pin_background_color"),
		AuthPagePinTextColor:        c.PostForm("auth_page_pin_text_color"),
		AuthPageTheme:               c.PostForm("auth_page_theme"),
		AuthPageCaptchaEnabled:      checkbox(c.PostForm("auth_page_captcha_enabled")),
		AuthPageCaptchaProvider:     c.PostForm("auth_page_captcha_provider"),
		AuthPageCaptchaSiteKey:      c.PostForm("auth_page_captcha_site_key"),
		AuthPageCaptchaSecretKey:    c.PostForm("auth_page_captcha_secret_key"),
		AuthPageBannerMode:          c.PostForm("auth_page_banner_mode"),
		AuthPageBanners:             parseBanners(c),
		AuthPageCopyCustom:          checkbox(c.PostForm("auth_page_copy_custom")),

		AuthPageLoginTitle:      c.PostForm("auth_page_login_title"),
		AuthPageLoginSubtitle:   c.PostForm("auth_page_login_subtitle"),
		AuthPageIdentifierLabel: c.PostForm("auth_page_identifier_label"),
		AuthPageSendButton:      c.PostForm("auth_page_send_button"),
		AuthPageVerifyTitle:     c.PostForm("auth_page_verify_title"),
		AuthPageVerifySubtitle:  c.PostForm("auth_page_verify_subtitle"),
		AuthPageVerifyButton:    c.PostForm("auth_page_verify_button"),
		AuthPageRetryText:       c.PostForm("auth_page_retry_text"),
		AuthPageBackText:        c.PostForm("auth_page_back_text"),
		AuthPageFooterText:      c.PostForm("auth_page_footer_text"),
	}
}

func parseBanners(c *gin.Context) []app.Banner {
	urls := c.PostFormArray("banner_url")
	texts := c.PostFormArray("banner_text")
	links := c.PostFormArray("banner_link")
	durations := c.PostFormArray("banner_duration")

	count := len(urls)
	for _, list := range [][]string{texts, links, durations} {
		if len(list) > count {
			count = len(list)
		}
	}

	banners := make([]app.Banner, 0, count)
	for i := 0; i < count; i++ {
		banner := app.Banner{
			URL:      strings.TrimSpace(valueAt(urls, i)),
			Text:     strings.TrimSpace(valueAt(texts, i)),
			Link:     strings.TrimSpace(valueAt(links, i)),
			Duration: atoiDefault(valueAt(durations, i), 6),
			// The order is the order of the rows in the form, which the admin
			// changes by dragging the banner cards.
			Order: i,
		}
		if banner.URL == "" && banner.Text == "" && banner.Link == "" {
			continue
		}
		banners = append(banners, banner)
	}
	return banners
}

func valueAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}

func pageParams(c *gin.Context, size int) (page, offset int) {
	page = atoiDefault(c.Query("page"), 1)
	if page < 1 {
		page = 1
	}
	return page, (page - 1) * size
}

func totalPages(total int64, size int) int {
	if total <= 0 || size <= 0 {
		return 0
	}
	pages := int((total + int64(size) - 1) / int64(size))
	if pages < 1 {
		pages = 1
	}
	return pages
}

func atoiDefault(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}

func checkbox(raw string) bool {
	return raw == "on" || raw == "true" || raw == "1"
}

func otpChannelLabel(channel string) string {
	switch channel {
	case app.OTPChannelEmail:
		return "ایمیل"
	case app.OTPChannelBoth:
		return "موبایل و ایمیل"
	default:
		return "پیامک"
	}
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// parseRange reads the dashboard range selector. "all" (or an empty default)
// selects every recorded session; the numeric values are day windows.
func parseRange(raw string) (days int, allTime bool) {
	switch strings.TrimSpace(raw) {
	case "all":
		return 0, true
	case "7":
		return 7, false
	case "14":
		return 14, false
	case "30":
		return 30, false
	case "90":
		return 90, false
	default:
		return 14, false
	}
}

func parseDateRange(c *gin.Context) (from, until *time.Time, rawFrom, rawTo, message string) {
	rawFrom = strings.TrimSpace(c.Query("from"))
	rawTo = strings.TrimSpace(c.Query("to"))
	from, until, err := adminpages.ParseJalaliRange(rawFrom, rawTo)
	if err != nil {
		return nil, nil, rawFrom, rawTo, err.Error()
	}
	return from, until, rawFrom, rawTo, ""
}

// chartStart picks the first day of the chart: the beginning of the selected
// range, or the first recorded day when the range is "all time".
func chartStart(since time.Time, counts []inbound.DailyCount, today time.Time) time.Time {
	if !since.IsZero() {
		return since
	}
	if len(counts) > 0 {
		return startOfDay(counts[0].Day)
	}
	return today
}

// daysBetween returns the number of daily points between the first recorded
// day and today, with at least one point so the chart never renders empty.
func daysBetween(counts []inbound.DailyCount, today time.Time) int {
	if len(counts) == 0 {
		return 1
	}
	start := startOfDay(counts[0].Day)
	days := int(today.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}
	return days
}

func toAppOptions(apps []app.App) []adminpages.AppOption {
	options := make([]adminpages.AppOption, len(apps))
	for i, a := range apps {
		options[i] = adminpages.AppOption{ID: a.ID.String(), Name: a.Name}
	}
	return options
}

func fillDaily(counts []inbound.DailyCount, since time.Time, days int) []inbound.DailyCount {
	byDay := make(map[string]int64, len(counts))
	for _, count := range counts {
		byDay[count.Day.Format("2006-01-02")] = count.Count
	}

	filled := make([]inbound.DailyCount, 0, days)
	for i := 0; i < days; i++ {
		day := since.AddDate(0, 0, i)
		filled = append(filled, inbound.DailyCount{
			Day:   day,
			Count: byDay[day.Format("2006-01-02")],
		})
	}
	return filled
}

// render wraps a page in the admin shell so every page gets the sidebar,
// top bar and the stylesheet link. Rendering a bare fragment would leave the
// page without CSS, which is what used to happen before the layout was wired.
func render(c *gin.Context, title, active string, component templ.Component) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := adminpages.AdminLayout(title, active, component).Render(c.Request.Context(), c.Writer); err != nil {
		redirectNotFound(c)
	}
}

func redirectNotFound(c *gin.Context) {
	c.Redirect(http.StatusFound, notFoundPath)
	c.Abort()
}

func redirectError(c *gin.Context, target, message string) {
	c.Redirect(http.StatusFound, target+"?error="+url.QueryEscape(message))
	c.Abort()
}
