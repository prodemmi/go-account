package pages

import (
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"aman/internal/adapters/inbound/http/pages"
	"aman/internal/ports/inbound"

	"github.com/a-h/templ"
)

// AppOption is one item of the app filter select.
type AppOption struct {
	ID   string
	Name string
}

// DashboardData is everything the admin dashboard renders.
type DashboardData struct {
	Apps        []AppOption
	AppID       string
	Days        int
	ChartDays   int
	AllTime     bool
	Stats       inbound.SessionStats
	Daily       []inbound.DailyCount
	GeneratedAt time.Time
	Error       string
}

func Dashboard(data DashboardData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := flash(data.Error, "").Render(ctx, w); err != nil {
			return err
		}

		var maxCount int64
		var peak inbound.DailyCount
		var sum int64
		for _, point := range data.Daily {
			sum += point.Count
			if point.Count > maxCount {
				maxCount = point.Count
				peak = point
			}
		}

		rangeLabel := rangeText(data)
		chartDays := data.ChartDays
		if chartDays <= 0 {
			chartDays = data.Days
		}
		if chartDays <= 0 {
			chartDays = len(data.Daily)
		}

		average := 0.0
		if chartDays > 0 {
			average = float64(sum) / float64(chartDays)
		}
		activeRate := percentage(data.Stats.Active, data.Stats.Total)
		todayTrend := int64(0)
		if len(data.Daily) >= 2 {
			todayTrend = data.Daily[len(data.Daily)-1].Count - data.Daily[len(data.Daily)-2].Count
		}

		peakLabel := "—"
		peakValue := int64(0)
		if maxCount > 0 {
			peakLabel = faNum(peak.Day.Format("2006-01-02"))
			peakValue = peak.Count
		}

		w.Write([]byte(`<section class="page-head dashboard-head">
	<div>
		<h2>داشبورد</h2>
		<p>نمای کلی ورودها، وضعیت نشست‌ها و روند روزانه.</p>
	</div>
	<span class="page-date">به‌روزرسانی: ` + faNum(FormatDateTime(data.GeneratedAt)) + `</span>
</section>

<section class="card dashboard-filter">
	<form method="GET" action="/admin" class="filter-form filters-form" data-auto-filter>
		<div class="field">
			<label for="app_id">
				برنامه
				<span class="help-tip" tabindex="0" aria-label="راهنمای فیلتر برنامه">?
					<span class="help-bubble">همه آمار را برای یک برنامه خاص ببینید. با «همه برنامه‌ها» داده‌های همه برنامه‌ها با هم جمع می‌شوند.</span>
				</span>
			</label>
			<select id="app_id" name="app_id">
				<option value="">همه برنامه‌ها</option>`))
		for _, option := range data.Apps {
			w.Write([]byte(`<option value="`))
			pages.WriteEscaped(w, option.ID)
			w.Write([]byte(`"`))
			if option.ID == data.AppID {
				w.Write([]byte(` selected`))
			}
			w.Write([]byte(`>`))
			pages.WriteEscaped(w, option.Name)
			w.Write([]byte(`</option>`))
		}
		w.Write([]byte(`</select>
		</div>
		<div class="field">
			<label for="days">
				بازه زمانی
				<span class="help-tip" tabindex="0" aria-label="راهنمای بازه زمانی">?
					<span class="help-bubble">نمودار و ترکیب وضعیت فقط نشست‌های ساخته‌شده در این بازه را نشان می‌دهند. «همه زمان‌ها» از اولین نشست ثبت‌شده حساب می‌کند.</span>
				</span>
			</label>
			<select id="days" name="days">`))
		writeRangeOption(w, "7", "۷ روز اخیر", data)
		writeRangeOption(w, "14", "۱۴ روز اخیر", data)
		writeRangeOption(w, "30", "۳۰ روز اخیر", data)
		writeRangeOption(w, "90", "۹۰ روز اخیر", data)
		writeRangeOption(w, "all", "همه زمان‌ها", data)
		w.Write([]byte(`</select>
		</div>
	</form>
</section>

<div class="stats-grid stats-grid-4 dashboard-stats">
	` + statCard("کل نشست‌ها", data.Stats.Total, "نشست‌های "+rangeLabel, "blue", iconTotal) + `
	` + statCard("فعال", data.Stats.Active, "قابل استفاده هم‌اکنون", "success", iconActive) + `
	` + statCard("معلق", data.Stats.Pending, "منتظر تبادل تیکت", "warning", iconPending) + `
	` + statCard("منقضی", data.Stats.Expired, "پایان‌یافته بدون ابطال", "danger", iconExpired) + `
	` + statCard("باطل‌شده", data.Stats.Revoked, "خروج یا ابطال دستی", "neutral", iconRevoked) + `
	` + statCard("امروز", data.Stats.Today, trendLabel(todayTrend), "violet", iconToday) + `
	` + statCard("کاربران یکتا", data.Stats.Unique, "شناسه‌های متمایز در بازه", "teal", iconUsers) + `
	` + statCard("نشست همیشگی", data.Stats.Infinite, "بدون انقضای عملی", "indigo", iconInfinite) + `
</div>

<div class="dashboard-grid dashboard-grid-wide">
	<section class="card chart-card">
		<div class="card-head">
			<div>
				<h3>روند ورود روزانه</h3>
				<span class="muted">` + rangeLabel + ` · ` + faNum(strconv.Itoa(chartDays)) + ` روز</span>
			</div>
			<div class="chart-badges">
				<span class="chart-badge">میانگین روزانه <strong>` + faNumFloat(average) + `</strong></span>
				<span class="chart-badge">اوج <strong>` + faNum(strconv.FormatInt(peakValue, 10)) + `</strong><small>` + peakLabel + `</small></span>
			</div>
		</div>`))

		if len(data.Daily) == 0 {
			w.Write([]byte(`<div class="filter-empty">در این بازه نشستی ثبت نشده است.</div>`))
		} else {
			if err := areaChart(w, data.Daily, maxCount); err != nil {
				return err
			}
		}

		w.Write([]byte(`</section>

	<section class="card breakdown-card">
		<div class="card-head">
			<div>
				<h3>ترکیب وضعیت نشست‌ها</h3>
				<span class="muted">` + rangeLabel + ` · مجموع ` + faNum(strconv.FormatInt(data.Stats.Total, 10)) + ` نشست</span>
			</div>
			<span class="chart-badge">نرخ فعال <strong>` + faNumFloat(activeRate) + `٪</strong></span>
		</div>
		<div class="breakdown-layout">
			` + donut(data.Stats) + `
			<div class="breakdown">`))

		breakdownRow(w, "فعال", data.Stats.Active, data.Stats.Total, "success")
		breakdownRow(w, "معلق", data.Stats.Pending, data.Stats.Total, "info")
		breakdownRow(w, "منقضی", data.Stats.Expired, data.Stats.Total, "warning")
		breakdownRow(w, "باطل‌شده", data.Stats.Revoked, data.Stats.Total, "neutral")

		w.Write([]byte(`</div>
		</div>
	</section>
</div>`))
		return nil
	})
}

// areaChart renders a dependency-free SVG area chart. Grid lines are fixed and
// the series is scaled to the largest day in the selected range.
func areaChart(w io.Writer, daily []inbound.DailyCount, maxCount int64) error {
	const width = 1000.0
	const height = 300.0
	const padTop = 18.0
	const padBottom = 34.0
	if len(daily) == 0 {
		return nil
	}
	if maxCount <= 0 {
		maxCount = 1
	}

	plotHeight := height - padTop - padBottom
	points := make([][2]float64, len(daily))
	for i, day := range daily {
		x := 0.0
		if len(daily) > 1 {
			x = float64(i) * width / float64(len(daily)-1)
		}
		y := padTop + plotHeight - (float64(day.Count)/float64(maxCount))*plotHeight
		points[i] = [2]float64{x, y}
	}

	var line strings.Builder
	for i, p := range points {
		if i == 0 {
			fmt.Fprintf(&line, "M %.2f %.2f", p[0], p[1])
		} else {
			fmt.Fprintf(&line, " L %.2f %.2f", p[0], p[1])
		}
	}
	area := fmt.Sprintf("%s L %.2f %.2f L %.2f %.2f Z", line.String(), width, height-padBottom, 0.0, height-padBottom)

	w.Write([]byte(`<div class="area-chart-wrap"><svg class="area-chart" viewBox="0 0 1000 300" preserveAspectRatio="none" role="img" aria-label="نمودار ورود روزانه">`))
	for _, pct := range []int{0, 25, 50, 75, 100} {
		y := padTop + plotHeight - (float64(pct)/100.0)*plotHeight
		w.Write([]byte(`<line class="chart-grid-line" x1="0" y1="` + fmt.Sprintf("%.2f", y) + `" x2="1000" y2="` + fmt.Sprintf("%.2f", y) + `" />`))
	}
	w.Write([]byte(`<defs><linearGradient id="area-fill" x1="0" y1="0" x2="0" y2="1"><stop offset="0%" stop-color="#2563eb" stop-opacity="0.30"/><stop offset="100%" stop-color="#2563eb" stop-opacity="0.02"/></linearGradient></defs>`))
	w.Write([]byte(`<path class="area-fill" d="` + area + `" fill="url(#area-fill)" />`))
	w.Write([]byte(`<path class="area-line" d="` + line.String() + `" vector-effect="non-scaling-stroke" />`))
	w.Write([]byte(`</svg></div>`))

	// Month/day labels stay sparse so long ranges never turn into noise.
	step := 1
	if len(daily) > 14 {
		step = int(math.Ceil(float64(len(daily)) / 10.0))
	}
	w.Write([]byte(`<div class="area-labels">`))
	for i, day := range daily {
		if i%step != 0 && i != len(daily)-1 {
			continue
		}
		w.Write([]byte(`<span>` + faNum(day.Day.Format("01/02")) + `</span>`))
	}
	w.Write([]byte(`</div>`))
	return nil
}

// donut renders the status split as a conic-gradient ring plus a center total.
func donut(stats inbound.SessionStats) string {
	total := stats.Total
	if total <= 0 {
		return `<div class="donut donut-empty"><span class="donut-center"><strong>۰</strong><small>نشست</small></span></div>`
	}
	segments := []struct {
		value int64
		color string
	}{
		{stats.Active, "#12b76a"},
		{stats.Pending, "#2563eb"},
		{stats.Expired, "#f79009"},
		{stats.Revoked, "#98a2b3"},
	}
	var gradient strings.Builder
	cursor := 0.0
	for _, segment := range segments {
		if segment.value <= 0 {
			continue
		}
		next := cursor + (float64(segment.value)/float64(total))*100
		fmt.Fprintf(&gradient, "%s %.2f%% %.2f%%, ", segment.color, cursor, next)
		cursor = next
	}
	if cursor < 100 {
		fmt.Fprintf(&gradient, "#eef2f7 %.2f%% 100%%, ", cursor)
	}
	style := strings.TrimSuffix(gradient.String(), ", ")
	return `<div class="donut" style="background: conic-gradient(` + style + `)"><span class="donut-center"><strong>` +
		faNum(strconv.FormatInt(total, 10)) + `</strong><small>نشست</small></span></div>`
}

func statCard(label string, value int64, note, tone, icon string) string {
	return `<article class="stat-card stat-card-` + tone + `">
		<div class="stat-card-head">
			<span class="stat-label">` + label + `</span>
			<span class="stat-icon stat-icon-` + tone + `">` + icon + `</span>
		</div>
		<strong class="stat-value">` + faNum(strconv.FormatInt(value, 10)) + `</strong>
		<small class="stat-note">` + note + `</small>
	</article>`
}

func trendLabel(trend int64) string {
	switch {
	case trend > 0:
		return "▲ " + faNum(strconv.FormatInt(trend, 10)) + " نسبت به دیروز"
	case trend < 0:
		return "▼ " + faNum(strconv.FormatInt(-trend, 10)) + " نسبت به دیروز"
	default:
		return "بدون تغییر نسبت به دیروز"
	}
}

func percentage(part, total int64) float64 {
	if total <= 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func faNumFloat(value float64) string {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return "۰"
	}
	text := fmt.Sprintf("%.1f", value)
	text = strings.TrimSuffix(text, ".0")
	return faNum(text)
}

func rangeText(data DashboardData) string {
	if data.AllTime {
		return "همه زمان‌ها"
	}
	switch data.Days {
	case 7:
		return "۷ روز اخیر"
	case 30:
		return "۳۰ روز اخیر"
	case 90:
		return "۹۰ روز اخیر"
	default:
		return "۱۴ روز اخیر"
	}
}

func writeRangeOption(w io.Writer, value, label string, data DashboardData) {
	selected := false
	if value == "all" {
		selected = data.AllTime
	} else if !data.AllTime {
		days, _ := strconv.Atoi(value)
		selected = data.Days == days
	}
	w.Write([]byte(`<option value="` + value + `"`))
	if selected {
		w.Write([]byte(` selected`))
	}
	w.Write([]byte(`>` + label + `</option>`))
}

func breakdownRow(w io.Writer, label string, value, total int64, tone string) {
	percent := 0
	if total > 0 {
		percent = int(math.Round(float64(value) / float64(total) * 100))
	}
	w.Write([]byte(`<div class="breakdown-row">
		<div class="breakdown-head">
			<span>` + label + `</span>
			<strong>` + faNum(strconv.FormatInt(value, 10)) + ` <small>` + faNum(strconv.Itoa(percent)) + `٪</small></strong>
		</div>
		<div class="progress"><span class="progress-fill tone-` + tone + `" style="width:` + fmt.Sprintf("%d", percent) + `%"></span></div>
	</div>`))
}

const (
	iconTotal    = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><rect x="3.5" y="4.5" width="17" height="15" rx="2.5"/><path d="M3.5 9.5h17M8 14h8"/></svg>`
	iconActive   = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3.5l7 2.8v5.2c0 4.1-2.9 7.7-7 9.5-4.1-1.8-7-5.4-7-9.5V6.3l7-2.8z"/><path d="M9.2 12.1l2 2 3.8-4"/></svg>`
	iconPending  = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="8.5"/><path d="M12 7.5V12l3 2"/></svg>`
	iconExpired  = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="8.5"/><path d="M12 8v4.5M12 16h.01M5.2 5.2l13.6 13.6"/></svg>`
	iconRevoked  = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="8.5"/><path d="M9 9l6 6M15 9l-6 6"/></svg>`
	iconToday    = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><rect x="3.5" y="5" width="17" height="15" rx="2.5"/><path d="M3.5 9.5h17M8 3v4M16 3v4"/><path d="M8 13h3v3H8z"/></svg>`
	iconUsers    = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="8" r="3"/><path d="M3.5 19c.6-3.2 2.7-5 5.5-5s4.9 1.8 5.5 5"/><path d="M15.5 6.2a3 3 0 0 1 0 5.6M17 14.4c2 .5 3.2 2 3.5 4.1"/></svg>`
	iconInfinite = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"><path d="M7 15.5c-2 0-3.5-1.6-3.5-3.5S5 8.5 7 8.5c3.4 0 6.6 7 10 7 2 0 3.5-1.6 3.5-3.5S19 8.5 17 8.5c-3.4 0-6.6 7-10 7z"/></svg>`
)
