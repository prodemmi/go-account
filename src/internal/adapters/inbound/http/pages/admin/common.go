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
	return m.BasePath + "?" + values.Encode()
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
