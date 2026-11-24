package pages

import (
	"context"
	"io"
	"net/url"
	"strconv"

	"aman/internal/adapters/inbound/http/pages"

	"github.com/a-h/templ"
)

// PageMeta describes server-side table filters and pagination.
type PageMeta struct {
	BasePath string
	Search   string
	Status   string
	AppID    string
	From     string
	To       string
	FilterError string
	Page     int
	Pages    int
	Total    int64
	PageSize int
}

// PageWindow returns up to five page numbers around the current page.
func (m PageMeta) PageWindow() []int {
	if m.Pages <= 0 {
		return nil
	}
	start := m.Page - 2
	if start < 1 {
		start = 1
	}
	end := start + 4
	if end > m.Pages {
		end = m.Pages
		start = end - 4
		if start < 1 {
			start = 1
		}
	}
	pages := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages
}

// Link builds a URL for a page preserving the active filters.
func (m PageMeta) Link(page int) string {
	values := url.Values{}
	values.Set("page", strconv.Itoa(page))
	if m.Search != "" {
		values.Set("q", m.Search)
	}
	if m.Status != "" {
		values.Set("status", m.Status)
	}
	if m.AppID != "" {
		values.Set("app_id", m.AppID)
	}
	if m.From != "" {
		values.Set("from", m.From)
	}
	if m.To != "" {
		values.Set("to", m.To)
	}
	return m.BasePath + "?" + values.Encode()
}

func dateRangeFields(w io.Writer, prefix, from, to string) {
	dateRangeField(w, prefix+"-from", "از تاریخ شمسی", from, "from")
	dateRangeField(w, prefix+"-to", "تا تاریخ شمسی", to, "to")
}

func dateRangeField(w io.Writer, id, label, value, name string) {
	w.Write([]byte(`<div class="field date-filter-field">
	<label for="`))
	pages.WriteEscaped(w, id)
	w.Write([]byte(`">`))
	pages.WriteEscaped(w, label)
	w.Write([]byte(`</label>
	<div class="jalali-date-field">
		<input type="text" id="`))
	pages.WriteEscaped(w, id)
	w.Write([]byte(`" name="`))
	pages.WriteEscaped(w, name)
	w.Write([]byte(`" value="`))
	pages.WriteEscaped(w, value)
	w.Write([]byte(`" placeholder="۱۴۰۵/۰۱/۰۱" inputmode="numeric" autocomplete="off" data-jalali-date />
		<button type="button" class="icon-action date-picker-toggle" data-date-open aria-label="انتخاب تاریخ" title="انتخاب تاریخ">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3.5" y="5" width="17" height="15.5" rx="2.5"/><path d="M7.5 3.5v3M16.5 3.5v3M3.5 9h17"/></svg>
		</button>
		<div class="jalali-calendar" data-jalali-calendar hidden></div>
	</div>
</div>`))
}

func copyButton(w io.Writer, value, label string) {
	w.Write([]byte(`<button type="button" class="copy-button" data-copy-text="`))
	pages.WriteEscaped(w, value)
	w.Write([]byte(`" aria-label="کپی `))
	pages.WriteEscaped(w, label)
	w.Write([]byte(`" title="کپی `))
	pages.WriteEscaped(w, label)
	w.Write([]byte(`">
		<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="8" y="8" width="12" height="12" rx="2"/><path d="M16 8V5a1.5 1.5 0 0 0-1.5-1.5h-9A1.5 1.5 0 0 0 4 5v9A1.5 1.5 0 0 0 5.5 15.5H8"/></svg>
	</button>`))
}

// CountLabel renders the total with a Persian unit.
func (m PageMeta) CountLabel(unit string) string {
	return faNum(strconv.FormatInt(m.Total, 10)) + " " + unit
}

func pagination(meta PageMeta) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if meta.Total == 0 {
			return nil
		}
		w.Write([]byte(`<nav class="pagination" aria-label="صفحه‌بندی">`))
		w.Write([]byte(`<span class="pagination-total">مجموع `))
		pages.WriteEscaped(w, faNum(strconv.FormatInt(meta.Total, 10)))
		w.Write([]byte(` مورد</span><div class="pagination-links">`))

		if meta.Page > 1 {
			w.Write([]byte(`<a class="page-link" href="`))
			pages.WriteEscaped(w, meta.Link(meta.Page-1))
			w.Write([]byte(`" aria-label="صفحه قبل">قبلی</a>`))
		} else {
			w.Write([]byte(`<span class="page-link is-disabled">قبلی</span>`))
		}

		for _, page := range meta.PageWindow() {
			class := "page-link"
			if page == meta.Page {
				class += " is-current"
			}
			w.Write([]byte(`<a class="` + class + `" href="`))
			pages.WriteEscaped(w, meta.Link(page))
			w.Write([]byte(`">`))
			pages.WriteEscaped(w, faNum(strconv.Itoa(page)))
			w.Write([]byte(`</a>`))
		}

		if meta.Page < meta.Pages {
			w.Write([]byte(`<a class="page-link" href="`))
			pages.WriteEscaped(w, meta.Link(meta.Page+1))
			w.Write([]byte(`" aria-label="صفحه بعد">بعدی</a>`))
		} else {
			w.Write([]byte(`<span class="page-link is-disabled">بعدی</span>`))
		}

		w.Write([]byte(`</div></nav>`))
		return nil
	})
}
