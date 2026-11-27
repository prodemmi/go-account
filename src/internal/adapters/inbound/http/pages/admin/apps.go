package pages

import (
	"context"
	"io"
	"strconv"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

func itostr(i int) string { return strconv.Itoa(i) }

// AppsList renders the admin apps list with a server-side search form and
// pagination.
func AppsList(apps []AppRow, meta PageMeta, errorMessage, successMessage string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		w.Write([]byte(`<section class="page-head">
	<div>
		<h2>برنامه‌ها</h2>
	</div>
	<a href="/admin/apps/new" class="btn btn-primary">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M12 5v14M5 12h14"/></svg>
		افزودن برنامه
	</a>
</section>`))

		if err := flash(errorMessage, successMessage).Render(ctx, w); err != nil {
			return err
		}

		w.Write([]byte(`<section class="card">
	<div class="card-head">
		<div>
			<h3>فهرست برنامه‌ها</h3>
			<span class="muted">`))
		pages.WriteEscaped(w, meta.CountLabel("برنامه"))
		w.Write([]byte(`</span>
		</div>
		<form method="GET" action="/admin/apps" class="filter-form filter-form-compact" data-auto-filter>
			<div class="search-box">
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><circle cx="11" cy="11" r="6.5"/><path d="M16 16l4 4"/></svg>
				<input type="search" name="q" value="`))
		pages.WriteEscaped(w, meta.Search)
		w.Write([]byte(`" placeholder="جستجو در نام، آدرس یا callback…" aria-label="جستجوی برنامه‌ها" />
			</div>`))
		dateRangeFields(w, "apps-date", meta.From, meta.To)
		w.Write([]byte(`<button type="submit" class="btn btn-ghost">جستجو</button>
			<button type="button" class="btn btn-ghost" data-clear-date-range>پیش‌فرض: همه زمان‌ها</button>
		</form>
	</div>`))
		if meta.FilterError != "" {
			w.Write([]byte(`<p class="date-filter-error" role="alert">`))
			pages.WriteEscaped(w, meta.FilterError)
			w.Write([]byte(`</p>`))
		}

		if len(apps) == 0 {
			w.Write([]byte(`<div class="empty-state">
		<span class="empty-icon">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6"><rect x="3.5" y="4.5" width="17" height="15" rx="2.5"/><path d="M3.5 9.5h17M8 4.5v5" stroke-linecap="round"/></svg>
		</span>
		<h4>`))
			if meta.Search != "" || meta.From != "" || meta.To != "" {
				w.Write([]byte(`نتیجه‌ای برای این جستجو پیدا نشد`))
			} else {
				w.Write([]byte(`هنوز برنامه‌ای ثبت نشده است`))
			}
			w.Write([]byte(`</h4>
		<p>`))
			if meta.Search != "" || meta.From != "" || meta.To != "" {
				w.Write([]byte(`عبارت جستجو را تغییر دهید یا همه برنامه‌ها را ببینید.`))
			} else {
				w.Write([]byte(`اولین برنامه را اضافه کنید تا جریان ورود و تنظیمات آن را مدیریت کنید.`))
			}
			w.Write([]byte(`</p>`))
			if meta.Search != "" || meta.From != "" || meta.To != "" {
				w.Write([]byte(`<a href="/admin/apps" class="btn btn-ghost">نمایش همه برنامه‌ها</a>`))
			} else {
				w.Write([]byte(`<a href="/admin/apps/new" class="btn btn-primary">افزودن اولین برنامه</a>`))
			}
			w.Write([]byte(`</div></section>`))
			return nil
		}

		w.Write([]byte(`<div class="table-wrap">
	<table class="table" id="apps-table">
		<thead>
			<tr>
				<th>برنامه</th>
				<th>روش ارسال کد</th>
				<th>اعتبار کد</th>
				<th>اعتبار تیکت</th>
				<th>مدت نشست</th>
				<th>صفحه ورود</th>
				<th>تاریخ ایجاد</th>
				<th class="th-actions">عملیات</th>
			</tr>
		</thead>
		<tbody>`))

		for _, a := range apps {
			w.Write([]byte(`<tr>
				<td>
					<div class="entity">
						<div class="entity-copy">
							<strong>`))
			pages.WriteEscaped(w, a.Name)
			w.Write([]byte(`</strong>
							<small class="app-id copy-id"><span class="ltr mono">`))
			pages.WriteEscaped(w, shortID(a.ID))
			w.Write([]byte(`</span>`))
			copyButton(w, a.ID, "شناسه برنامه")
			w.Write([]byte(`</small>
						</div>
					</div>
				</td>
				<td><span class="badge badge-info">`))
			pages.WriteEscaped(w, a.OTPChannel)
			w.Write([]byte(`</span></td>
				<td class="nowrap">`))
			pages.WriteEscaped(w, humanSeconds(a.OTPExpire))
			w.Write([]byte(`</td>
				<td class="nowrap">`))
			pages.WriteEscaped(w, humanSeconds(a.OTTExpire))
			w.Write([]byte(`</td>
				<td class="nowrap">`))
			if a.InfiniteSession {
				w.Write([]byte(`<span class="badge badge-violet">همیشگی</span>`))
			} else {
				pages.WriteEscaped(w, humanSeconds(a.SessionTimeout))
			}
			w.Write([]byte(`</td>
				<td>`))
			if a.AuthPageEnabled {
				w.Write([]byte(`<span class="badge badge-success">سفارشی</span>`))
			} else {
				w.Write([]byte(`<span class="badge badge-neutral">پیش‌فرض</span>`))
			}
			w.Write([]byte(`</td>
				<td class="nowrap">`))
			pages.WriteEscaped(w, a.CreatedAt)
			w.Write([]byte(`</td>
				<td>
					<div class="row-actions">
						<a href="/admin/apps/`))
			pages.WriteEscaped(w, a.ID)
			w.Write([]byte(`" class="icon-action" title="ویرایش برنامه" aria-label="ویرایش برنامه">
							<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4 20h4.2L19 9.2a2.1 2.1 0 0 0-3-3L5.2 17 4 20z"/><path d="M14.5 7.5l2 2"/></svg>
						</a>
						<form method="POST" action="/admin/apps/`))
			pages.WriteEscaped(w, a.ID)
			w.Write([]byte(`/delete" data-confirm="این برنامه و تنظیمات آن حذف شود؟ این عملیات قابل بازگشت نیست.">
							<button type="submit" class="icon-action icon-action-danger" title="حذف برنامه" aria-label="حذف برنامه">
								<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M5 7h14M10 7V5h4v2M7 7l.7 12h8.6L17 7"/><path d="M10.5 10.5v5M13.5 10.5v5"/></svg>
							</button>
						</form>
					</div>
				</td>
			</tr>`))
		}

		w.Write([]byte(`		</tbody>
	</table>
</div>`))
		if err := pagination(meta).Render(ctx, w); err != nil {
			return err
		}
		w.Write([]byte(`</section>`))
		return nil
	})
}

type AppRow struct {
	ID              string
	Name            string
	Url             string
	CallbackURL     string
	OTPChannel      string
	OTPExpire       int
	OTTExpire       int
	SessionTimeout  int
	InfiniteSession bool
	AuthPageEnabled bool
	CreatedAt       string
}
