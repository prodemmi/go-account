package pages

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"aman/internal/ports/inbound"
)

func TestWeekdayCountsUseSaturdayFirstAndAggregateDailyCounts(t *testing.T) {
	daily := []inbound.DailyCount{
		{Day: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), Count: 2}, // Monday
		{Day: time.Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC), Count: 3}, // Saturday
		{Day: time.Date(2024, time.January, 7, 0, 0, 0, 0, time.UTC), Count: 4}, // Sunday
		{Day: time.Date(2024, time.January, 8, 0, 0, 0, 0, time.UTC), Count: 5}, // Monday
	}
	want := [7]int64{3, 4, 7, 0, 0, 0, 0}
	if got := weekdayCounts(daily); got != want {
		t.Fatalf("weekdayCounts() = %v, want %v", got, want)
	}
}

func TestDashboardRendersSelectedStatusAndWeekdayChart(t *testing.T) {
	var output bytes.Buffer
	data := DashboardData{
		Status:      "expired",
		Days:        7,
		ChartDays:   7,
		GeneratedAt: time.Date(2024, time.January, 8, 12, 0, 0, 0, time.UTC),
		Daily: []inbound.DailyCount{
			{Day: time.Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC), Count: 3},
		},
	}
	if err := Dashboard(data).Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<select id="dashboard-status" name="status">`,
		`<option value="expired" selected>منقضی</option>`,
		`تغییر روز آخر`,
		`نشست‌ها به تفکیک روز هفته`,
		`<span class="weekday-count">۳</span>`,
	} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("dashboard is missing %q", want)
		}
	}
}
