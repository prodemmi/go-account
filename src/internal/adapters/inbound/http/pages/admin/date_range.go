package pages

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseJalaliRange parses optional Jalali dates into a half-open Gregorian
// range. The end is exclusive so it includes the whole selected final day.
func ParseJalaliRange(from, to string) (start, end *time.Time, err error) {
	if strings.TrimSpace(from) != "" {
		value, parseErr := parseJalaliDate(from)
		if parseErr != nil {
			return nil, nil, fmt.Errorf("تاریخ شروع معتبر نیست؛ قالب سال/ماه/روز شمسی را وارد کنید")
		}
		start = &value
	}
	if strings.TrimSpace(to) != "" {
		value, parseErr := parseJalaliDate(to)
		if parseErr != nil {
			return nil, nil, fmt.Errorf("تاریخ پایان معتبر نیست؛ قالب سال/ماه/روز شمسی را وارد کنید")
		}
		value = value.AddDate(0, 0, 1)
		end = &value
	}
	if start != nil && end != nil && !start.Before(*end) {
		return nil, nil, fmt.Errorf("تاریخ پایان باید برابر یا بعد از تاریخ شروع باشد")
	}
	return start, end, nil
}

func parseJalaliDate(raw string) (time.Time, error) {
	var ascii strings.Builder
	for _, r := range strings.TrimSpace(raw) {
		switch {
		case r >= '۰' && r <= '۹':
			ascii.WriteByte(byte('0' + r - '۰'))
		case r >= '٠' && r <= '٩':
			ascii.WriteByte(byte('0' + r - '٠'))
		default:
			ascii.WriteRune(r)
		}
	}
	parts := strings.FieldsFunc(ascii.String(), func(r rune) bool { return r == '/' || r == '-' })
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid Jalali date")
	}
	jy, err := strconv.Atoi(parts[0])
	if err != nil || jy < 1 || jy > 9998 {
		return time.Time{}, fmt.Errorf("invalid Jalali year")
	}
	jm, err := strconv.Atoi(parts[1])
	if err != nil || jm < 1 || jm > 12 {
		return time.Time{}, fmt.Errorf("invalid Jalali month")
	}
	jd, err := strconv.Atoi(parts[2])
	if err != nil || jd < 1 || jd > 31 {
		return time.Time{}, fmt.Errorf("invalid Jalali day")
	}

	// The Jalali year begins around March 20/21. Searching that Gregorian year
	// with the existing conversion keeps this parser aligned with FormatDateTime.
	first := time.Date(jy+621, time.March, 19, 0, 0, 0, 0, time.Local)
	for offset := 0; offset < 370; offset++ {
		candidate := first.AddDate(0, 0, offset)
		y, m, d := candidate.Date()
		cy, cm, cd := toJalali(y, int(m), d)
		if cy == jy && cm == jm && cd == jd {
			return candidate, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid Jalali date")
}
