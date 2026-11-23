package pages

import (
	"testing"
	"time"
)

func TestFormatDateTimeUsesJalaliCalendar(t *testing.T) {
	at := time.Date(2026, time.November, 17, 23, 31, 0, 0, time.Local)
	got := FormatDateTime(at)
	want := "۱۴۰۵/۰۸/۲۶ · ۲۳:۳۱"
	if got != want {
		t.Fatalf("FormatDateTime() = %q, want %q", got, want)
	}
}

func TestToJalaliNowruz(t *testing.T) {
	jy, jm, jd := toJalali(2024, 3, 20)
	if jy != 1403 || jm != 1 || jd != 1 {
		t.Fatalf("toJalali(2024-03-20) = %d/%d/%d, want 1403/1/1", jy, jm, jd)
	}
}
