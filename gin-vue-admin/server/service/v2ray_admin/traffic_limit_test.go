package v2ray_admin

import (
	"testing"
	"time"
)

func TestMonthlyTrafficLimitStartCreatedAtUsesCurrentMonthStart(t *testing.T) {
	now := mustShanghaiDate(t, 2026, 6, 22)

	got := MonthlyTrafficLimitStartCreatedAt(now)

	if got != 20260601 {
		t.Fatalf("start created_at = %d, want 20260601", got)
	}
}

func TestMonthlyTrafficLimitStartCreatedAtDoesNotUseServerResetDay(t *testing.T) {
	now := mustShanghaiDate(t, 2026, 6, 19)

	got := MonthlyTrafficLimitStartCreatedAt(now)

	if got != 20260601 {
		t.Fatalf("start created_at = %d, want 20260601", got)
	}
}

func TestMonthlyTrafficLimitStartCreatedAtUsesFirstDayWhenServerResetWouldClamp(t *testing.T) {
	now := mustShanghaiDate(t, 2026, 2, 28)

	got := MonthlyTrafficLimitStartCreatedAt(now)

	if got != 20260201 {
		t.Fatalf("start created_at = %d, want 20260201", got)
	}
}

func TestIsMonthlyTrafficLimitResetDayUsesNaturalMonthStart(t *testing.T) {
	if !IsMonthlyTrafficLimitResetDay(mustShanghaiDate(t, 2026, 6, 1)) {
		t.Fatal("June 1 should be monthly traffic limit reset day")
	}
	if IsMonthlyTrafficLimitResetDay(mustShanghaiDate(t, 2026, 6, 20)) {
		t.Fatal("June 20 should not be monthly traffic limit reset day")
	}
}

func TestIsTrafficResetDayUsesClampedServerResetDay(t *testing.T) {
	if !IsTrafficResetDay(mustShanghaiDate(t, 2026, 6, 20), 20) {
		t.Fatal("June 20 should be reset day when reset day is 20")
	}
	if IsTrafficResetDay(mustShanghaiDate(t, 2026, 6, 21), 20) {
		t.Fatal("June 21 should not be reset day when reset day is 20")
	}
	if !IsTrafficResetDay(mustShanghaiDate(t, 2026, 2, 28), 31) {
		t.Fatal("Feb 28 should be reset day when reset day 31 is clamped to month end")
	}
}

func TestShouldLimitTraffic(t *testing.T) {
	if ShouldLimitTraffic(10, -1) {
		t.Fatal("negative limit should mean unlimited")
	}
	if !ShouldLimitTraffic(0, 0) {
		t.Fatal("zero limit should mean always limited")
	}
	if ShouldLimitTraffic(1024*1024*1024-1, 1) {
		t.Fatal("traffic below limit should not be limited")
	}
	if !ShouldLimitTraffic(1024*1024*1024, 1) {
		t.Fatal("traffic equal to limit should be limited")
	}
}

func mustShanghaiDate(t *testing.T, year int, month time.Month, day int) time.Time {
	t.Helper()
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	return time.Date(year, month, day, 12, 0, 0, 0, location)
}
