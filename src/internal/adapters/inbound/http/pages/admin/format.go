package pages

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// faNum converts ASCII digits to Persian digits for display.
func faNum(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune('۰' + (r - '0'))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// humanSeconds renders a duration in the largest sensible Persian unit.
func humanSeconds(seconds int) string {
	if seconds <= 0 {
		return "—"
	}
	switch {
	case seconds%86400 == 0 && seconds >= 86400:
		return faNum(strconv.Itoa(seconds/86400)) + " روز"
	case seconds%3600 == 0 && seconds >= 3600:
		return faNum(strconv.Itoa(seconds/3600)) + " ساعت"
	case seconds%60 == 0 && seconds >= 60:
		return faNum(strconv.Itoa(seconds/60)) + " دقیقه"
	default:
		return faNum(strconv.Itoa(seconds)) + " ثانیه"
	}
}

// shortID keeps long identifiers readable inside dense tables.
func shortID(id string) string {
	if len(id) <= 10 {
		return id
	}
	return id[:8] + "…"
}

// firstRune returns the first character of a name, used by avatar chips.
func firstRune(s string) string {
	s = strings.TrimSpace(s)
	for _, r := range s {
		return strings.ToUpper(string(r))
	}
	return "?"
}

// FormatDateTime renders a timestamp as a Jalali date and time with Persian
// digits, e.g. ۱۴۰۵/۰۸/۲۶ · ۲۳:۳۱. Persian admins read Jalali dates, not the
// Gregorian timestamps stored in the database.
func FormatDateTime(t time.Time) string {
	y, m, d := t.Date()
	jy, jm, jd := toJalali(y, int(m), d)
	raw := fmt.Sprintf("%04d/%02d/%02d · %02d:%02d", jy, jm, jd, t.Hour(), t.Minute())
	return faNum(raw)
}

// toJalali converts a Gregorian date to the Jalali calendar. It is the
// standard arithmetic conversion and needs no external dependency.
func toJalali(gy, gm, gd int) (jy, jm, jd int) {
	monthDays := [...]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}

	gy2 := gy
	if gm > 2 {
		gy2++
	}
	days := 355666 + (365 * gy) + ((gy2 + 3) / 4) - ((gy2 + 99) / 100) + ((gy2 + 399) / 400) + gd + monthDays[gm-1]

	jy = -1595 + (33 * (days / 12053))
	days %= 12053
	jy += 4 * (days / 1461)
	days %= 1461
	if days > 365 {
		jy += (days - 1) / 365
		days = (days - 1) % 365
	}
	if days < 186 {
		jm = 1 + days/31
		jd = 1 + days%31
	} else {
		jm = 7 + (days-186)/30
		jd = 1 + (days-186)%30
	}
	return jy, jm, jd
}
