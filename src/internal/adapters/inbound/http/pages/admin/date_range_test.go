package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestParseJalaliRangeUsesPersianDigitsAndIncludesEndDay(t *testing.T) {
	start, end, err := ParseJalaliRange("۱۴۰۳/۰۱/۰۱", "١٤٠٣-٠١-٠٢")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := start.Format("2006-01-02"), "2024-03-20"; got != want {
		t.Fatalf("start = %s, want %s", got, want)
	}
	if got, want := end.Format("2006-01-02"), "2024-03-22"; got != want {
		t.Fatalf("exclusive end = %s, want %s", got, want)
	}
	if start.Hour() != 0 || end.Hour() != 0 || start.Location() != time.Local || end.Location() != time.Local {
		t.Fatalf("range should use local day boundaries: start=%v end=%v", start, end)
	}
}

func TestParseJalaliRangeRejectsInvalidAndReversedRanges(t *testing.T) {
	for _, tc := range []struct {
		name, from, to string
	}{
		{name: "invalid day", from: "۱۴۰۳/۰۷/۳۱"},
		{name: "reversed", from: "۱۴۰۳/۰۱/۰۳", to: "۱۴۰۳/۰۱/۰۲"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := ParseJalaliRange(tc.from, tc.to); err == nil {
				t.Fatal("ParseJalaliRange() accepted an invalid range")
			}
		})
	}
}

func TestDateRangeControlsAppearOnAdminTablesAndDashboard(t *testing.T) {
	components := map[string]templ.Component{
		"apps table":     AppsList(nil, PageMeta{BasePath: "/admin/apps"}, "", ""),
		"sessions table": SessionsList(nil, PageMeta{BasePath: "/admin/sessions"}, nil, "", ""),
		"dashboard":      Dashboard(DashboardData{}),
	}
	for name, component := range components {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			if err := component.Render(context.Background(), &output); err != nil {
				t.Fatal(err)
			}
			html := output.String()
			if got := strings.Count(html, `data-jalali-date`); got != 2 {
				t.Fatalf("rendered %d Jalali date fields, want 2", got)
			}
			if !strings.Contains(html, `data-clear-date-range`) || !strings.Contains(html, `پیش‌فرض: همه زمان‌ها`) {
				t.Fatal("all-time reset control is missing")
			}
		})
	}
}
