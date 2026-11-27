package pages

import (
	"context"
	"io"
	"strings"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

// SessionsList renders the admin sessions list with server-side filters and
// pagination. It intentionally has no "online" badge: status is derived from
// the persisted session lifecycle (pending, active, expired, revoked).
func SessionsList(sessions []SessionRow, meta PageMeta, apps []AppOption, errorMessage, successMessage string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<section class="page-head">
	<div>
		<h2>نشست‌ها</h2>
	</div>
</section>`))

		if err := flash(errorMessage, successMessage).Render(ctx, w); err != nil {
			return err
		}

		// Filters use the same compact card layout as the dashboard. They
		// apply themselves as soon as a value changes.
		w.Write([]byte(`<section class="card filters-card">
	<form method="GET" action="/admin/sessions" class="filter-form filters-form" data-auto-filter>
		<div class="field filter-search">
			<label for="sessions-q">جستجو</label>
			<div class="search-box">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><circle cx="11" cy="11" r="6.5"/><path d="M16 16l4 4"/></svg>
				<input type="search" id="sessions-q" name="q" value="`))
		pages.WriteEscaped(w, meta.Search)
		w.Write([]byte(`" placeholder="کاربر، IP، شهر یا مرورگر…" />
			</div>
		</div>
		<div class="field">
			<label for="sessions-app">برنامه</label>
			<select id="sessions-app" name="app_id">
				<option value="">همه برنامه‌ها</option>`))
		for _, option := range apps {
			w.Write([]byte(`<option value="`))
			pages.WriteEscaped(w, option.ID)
			w.Write([]byte(`"`))
			if option.ID == meta.AppID {
				w.Write([]byte(` selected`))
			}
			w.Write([]byte(`>`))
			pages.WriteEscaped(w, option.Name)
			w.Write([]byte(`</option>`))
		}
		w.Write([]byte(`</select>
		</div>
		<div class="field">
			<label for="sessions-status">وضعیت</label>
			<select id="sessions-status" name="status">
				<option value="">همه وضعیت‌ها</option>`))
		writeStatusOption(w, "active", "فعال", meta.Status)
		writeStatusOption(w, "pending", "معلق", meta.Status)
		writeStatusOption(w, "expired", "منقضی", meta.Status)
		writeStatusOption(w, "revoked", "باطل‌شده", meta.Status)
		w.Write([]byte(`</select>
		</div>`))
		dateRangeFields(w, "sessions-date", meta.From, meta.To)
		w.Write([]byte(`<button type="button" class="btn btn-ghost" data-clear-date-range>پیش‌فرض: همه زمان‌ها</button>
	</form>`))
		if meta.FilterError != "" {
			w.Write([]byte(`<p class="date-filter-error" role="alert">`))
			pages.WriteEscaped(w, meta.FilterError)
			w.Write([]byte(`</p>`))
		}
		w.Write([]byte(`
</section>

<section class="card">
	<div class="card-head">
		<div>
			<h3>فهرست نشست‌ها</h3>
			<span class="muted">`))
		pages.WriteEscaped(w, meta.CountLabel("نشست"))
		w.Write([]byte(`</span>
		</div>
	</div>`))

		if len(sessions) == 0 {
			w.Write([]byte(`<div class="empty-state">
		<span class="empty-icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round"><rect x="3.5" y="5" width="17" height="14" rx="2.5"/><path d="M8 12h8M8 9.5h4"/></svg>
		</span>
		<h4>`))
			if meta.Search != "" || meta.Status != "" || meta.AppID != "" || meta.From != "" || meta.To != "" {
				w.Write([]byte(`نتیجه‌ای برای این فیلترها پیدا نشد`))
			} else {
				w.Write([]byte(`هنوز نشستی ثبت نشده است`))
			}
			w.Write([]byte(`</h4>
		<p>`))
			if meta.Search != "" || meta.Status != "" || meta.AppID != "" || meta.From != "" || meta.To != "" {
				w.Write([]byte(`فیلترها را تغییر دهید تا نتیجه دیگری نمایش داده شود.`))
			} else {
				w.Write([]byte(`پس از اولین ورود کاربران از طریق برنامه‌ها، نشست‌ها اینجا نمایش داده می‌شوند.`))
			}
			w.Write([]byte(`</p></div></section>`))
			return nil
		}

		w.Write([]byte(`<div class="table-wrap">
	<table class="table sessions-table" id="sessions-table">
		<thead>
			<tr>
				<th>شناسه کاربر</th>
				<th>برنامه</th>
				<th>IP</th>
				<th>دستگاه</th>
				<th>مرورگر</th>
				<th>وضعیت</th>
				<th class="th-actions">عملیات</th>
			</tr>
		</thead>
		<tbody>`))

		for i, s := range sessions {
			drawerID := "session-detail-" + itostr(i)

			w.Write([]byte(`<tr>
				<td><strong class="ltr table-strong">`))
			pages.WriteEscaped(w, s.Identifier)
			w.Write([]byte(`</strong><small class="copy-id"><span class="ltr mono">`))
			pages.WriteEscaped(w, shortID(s.ID))
			w.Write([]byte(`</span>`))
			copyButton(w, s.ID, "شناسه نشست")
			w.Write([]byte(`</small></td>
				<td><strong class="table-strong">`))
			pages.WriteEscaped(w, s.AppName)
			w.Write([]byte(`</strong></td>
				<td class="ltr mono nowrap">`))
			if s.IP != "" {
				pages.WriteEscaped(w, s.IP)
			} else {
				w.Write([]byte(`—`))
			}
			w.Write([]byte(`</td>
				<td>`))
			if s.UserAgent != "" || s.Device != "" {
				pages.WriteEscaped(w, humanDeviceName(s.UserAgent, s.Device))
			} else {
				w.Write([]byte(`—`))
			}
			w.Write([]byte(`</td>
				<td>`))
			pages.WriteEscaped(w, browserLabel(s.UserAgent))
			w.Write([]byte(`</td>
				<td>`))
			writeStatusBadge(w, s.Status)
			w.Write([]byte(`</td>
				<td>
					<div class="row-actions">
						<button type="button" class="icon-action" title="جزئیات نشست" aria-label="جزئیات نشست" data-drawer-open="`))
			pages.WriteEscaped(w, drawerID)
			w.Write([]byte(`">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="8.5"/><path d="M12 11v5M12 8h.01"/></svg>
						</button>
						<form method="POST" action="/admin/sessions/`))
			pages.WriteEscaped(w, s.ID)
			w.Write([]byte(`/delete" data-confirm="این نشست ابطال شود؟ کاربر باید دوباره وارد شود.">
							<button type="submit" class="icon-action icon-action-danger" title="ابطال نشست" aria-label="ابطال نشست">
								<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3.5l7 2.8v5.2c0 4.1-2.9 7.7-7 9.5-4.1-1.8-7-5.4-7-9.5V6.3l7-2.8z"/><path d="M9.5 9.5l5 5M14.5 9.5l-5 5"/></svg>
							</button>
						</form>
					</div>
				</td>
			</tr>`))
		}

		w.Write([]byte(`		</tbody>
	</table>
</div>`))
		// Drawers live outside the table: a dialog is not valid table content.
		for i, s := range sessions {
			drawerID := "session-detail-" + itostr(i)
			if err := sessionDrawer(drawerID, s).Render(ctx, w); err != nil {
				return err
			}
		}
		if err := pagination(meta).Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`</section>`))
		return nil
	})
}

func sessionDrawer(id string, s SessionRow) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		w.Write([]byte(`<dialog class="drawer" aria-label="جزئیات نشست" id="`))
		pages.WriteEscaped(w, id)
		w.Write([]byte(`">
	<div class="drawer-head">
		<div>
			<h3>جزئیات نشست</h3>
			<span class="muted small">`))
		pages.WriteEscaped(w, s.Identifier)
		w.Write([]byte(`</span>
		</div>
		<button type="button" class="icon-btn" data-drawer-close aria-label="بستن">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M6 6l12 12M18 6L6 18"/></svg>
		</button>
	</div>
	<div class="drawer-body">
		<div class="detail-status">`))
		writeStatusBadge(w, s.Status)
		w.Write([]byte(`</div>
		<dl class="detail-list">
			<div class="detail-row"><dt>شناسه نشست</dt><dd class="copy-value ltr mono"><span>`))
		pages.WriteEscaped(w, s.ID)
		w.Write([]byte(`</span>`))
		copyButton(w, s.ID, "شناسه نشست")
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>شناسه برنامه</dt><dd class="copy-value ltr mono"><span>`))
		pages.WriteEscaped(w, s.AppID)
		w.Write([]byte(`</span>`))
		copyButton(w, s.AppID, "شناسه برنامه")
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>نام برنامه</dt><dd>`))
		pages.WriteEscaped(w, s.AppName)
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>شناسه کاربر</dt><dd class="copy-value ltr"><span>`))
		pages.WriteEscaped(w, s.Identifier)
		w.Write([]byte(`</span>`))
		copyButton(w, s.Identifier, "شناسه کاربر")
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>IP</dt><dd class="ltr mono">`))
		if s.IP != "" {
			pages.WriteEscaped(w, s.IP)
		} else {
			w.Write([]byte(`—`))
		}
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>موقعیت</dt><dd>`))
		location := joinLocation(s.City, s.Country)
		if s.Region != "" {
			if location != "" {
				location += "، "
			}
			location += s.Region
		}
		if location == "" {
			location = "—"
		}
		pages.WriteEscaped(w, location)
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>دستگاه</dt><dd>`))
		if s.UserAgent != "" || s.Device != "" {
			pages.WriteEscaped(w, humanDeviceName(s.UserAgent, s.Device))
		} else {
			w.Write([]byte(`—`))
		}
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>مرورگر</dt><dd>`))
		pages.WriteEscaped(w, browserLabel(s.UserAgent))
		w.Write([]byte(`</dd></div>
			<div class="detail-row detail-row-block"><dt>User-Agent</dt><dd class="ltr mono small">`))
		if s.UserAgent != "" {
			pages.WriteEscaped(w, s.UserAgent)
		} else {
			w.Write([]byte(`—`))
		}
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>ساخت</dt><dd>`))
		pages.WriteEscaped(w, s.CreatedAt)
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>آخرین ورود</dt><dd>`))
		pages.WriteEscaped(w, s.LastLoginAt)
		w.Write([]byte(`</dd></div>
			<div class="detail-row"><dt>انقضا</dt><dd>`))
		if s.Infinite {
			w.Write([]byte(`بدون انقضا (نشست همیشگی)`))
		} else {
			pages.WriteEscaped(w, s.ExpiresAt)
		}
		w.Write([]byte(`</dd></div>
		</dl>
		<form method="POST" action="/admin/sessions/`))
		pages.WriteEscaped(w, s.ID)
		w.Write([]byte(`/delete" data-confirm="این نشست ابطال شود؟ کاربر باید دوباره وارد شود." class="drawer-actions">
			<button type="submit" class="btn btn-danger btn-block">ابطال این نشست</button>
		</form>
	</div>
</dialog>`))
		return nil
	})
}

func writeStatusBadge(w io.Writer, status string) {
	switch status {
	case "active":
		w.Write([]byte(`<span class="badge badge-success">فعال</span>`))
	case "pending":
		w.Write([]byte(`<span class="badge badge-info">معلق</span>`))
	case "expired":
		w.Write([]byte(`<span class="badge badge-warning">منقضی</span>`))
	default:
		w.Write([]byte(`<span class="badge badge-neutral">باطل‌شده</span>`))
	}
}

func writeStatusOption(w io.Writer, value, label, selected string) {
	w.Write([]byte(`<option value="` + value + `"`))
	if value == selected {
		w.Write([]byte(` selected`))
	}
	w.Write([]byte(`>` + label + `</option>`))
}

func joinLocation(city, country string) string {
	switch {
	case city != "" && country != "":
		return city + "، " + country
	case city != "":
		return city
	default:
		return country
	}
}

func deviceLabel(device string) string {
	switch device {
	case "mobile":
		return "موبایل"
	case "tablet":
		return "تبلت"
	case "desktop":
		return "دسکتاپ"
	case "bot":
		return "ربات"
	default:
		return device
	}
}

func humanDeviceName(userAgent, fallback string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "ipad"):
		return "iPad"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipod"):
		return "iPhone"
	}

	browser := browserLabel(userAgent)
	operatingSystem := ""
	switch {
	case strings.Contains(ua, "cros"):
		operatingSystem = "ChromeOS"
	case strings.Contains(ua, "windows"):
		operatingSystem = "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		operatingSystem = "macOS"
	case strings.Contains(ua, "android"):
		operatingSystem = "Android"
	case strings.Contains(ua, "linux"):
		operatingSystem = "Linux"
	case strings.Contains(ua, "ios"):
		operatingSystem = "iOS"
	}
	if browser != "نامشخص" && browser != "—" && operatingSystem != "" {
		return browser + " on " + operatingSystem
	}
	if operatingSystem != "" {
		return operatingSystem
	}
	if browser != "نامشخص" && browser != "—" {
		return browser
	}
	if fallback != "" {
		return deviceLabel(fallback)
	}
	return "نامشخص"
}

// browserLabel extracts the human browser family from a User-Agent string. It
// is only used for display; the full agent stays available in the drawer.
func browserLabel(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case ua == "":
		return "—"
	case strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/"):
		return "Edge"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		return "Opera"
	case strings.Contains(ua, "samsungbrowser"):
		return "Samsung Internet"
	case strings.Contains(ua, "firefox") || strings.Contains(ua, "fxios"):
		return "Firefox"
	case strings.Contains(ua, "chrome") || strings.Contains(ua, "crios"):
		return "Chrome"
	case strings.Contains(ua, "safari"):
		return "Safari"
	case strings.Contains(ua, "curl"):
		return "curl"
	case strings.Contains(ua, "wget"):
		return "wget"
	case strings.Contains(ua, "bot"):
		return "ربات"
	default:
		return "نامشخص"
	}
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit] + "…"
}

type SessionRow struct {
	ID          string
	AppID       string
	AppName     string
	Identifier  string
	IP          string
	Country     string
	City        string
	Region      string
	Device      string
	UserAgent   string
	LastLoginAt string
	CreatedAt   string
	ExpiresAt   string
	Status      string
	Infinite    bool
}
