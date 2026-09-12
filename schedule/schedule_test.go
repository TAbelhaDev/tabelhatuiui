package schedule

import (
	"strings"
	"testing"
	"time"
)

func TestOnCalendarDaily(t *testing.T) {
	s := Schedule{Kind: KindDaily, Hour: 21, Minute: 0}
	got := s.OnCalendar(time.Now())
	want := "*-*-* 21:00:00"
	if got != want {
		t.Errorf("daily: got %q, want %q", got, want)
	}
}

func TestOnCalendarWeekly(t *testing.T) {
	s := Schedule{Kind: KindWeekly, Hour: 9, Minute: 0, Weekdays: []time.Weekday{time.Monday, time.Friday}}
	got := s.OnCalendar(time.Now())
	want := "Mon,Fri *-*-* 09:00:00"
	if got != want {
		t.Errorf("weekly: got %q, want %q", got, want)
	}
}

func TestOnCalendarMonthly(t *testing.T) {
	s := Schedule{Kind: KindMonthly, Hour: 8, Minute: 30, DayOfMonth: 15}
	got := s.OnCalendar(time.Now())
	want := "*-*-15 08:30:00"
	if got != want {
		t.Errorf("monthly: got %q, want %q", got, want)
	}
}

func TestOnCalendarOneshot(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.Local)
	s := Schedule{Kind: KindOneshot, Hour: 21, Minute: 0, DOM: 15, Month: 9}
	got := s.OnCalendar(now)
	want := "2026-09-15 21:00:00"
	if got != want {
		t.Errorf("oneshot: got %q, want %q", got, want)
	}
}

func TestOnCalendarOneshotRollover(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.Local)
	s := Schedule{Kind: KindOneshot, Hour: 21, Minute: 0, DOM: 1, Month: 8}
	got := s.OnCalendar(now)
	want := "2027-08-01 21:00:00"
	if got != want {
		t.Errorf("oneshot rollover: got %q, want %q", got, want)
	}
}

func TestOnCalendarCycle(t *testing.T) {
	now := time.Date(2026, 9, 12, 10, 0, 0, 0, time.Local)
	s := Schedule{Kind: KindCycle, Hour: 14, Minute: 0, DOM: 12, Month: 9}
	got := s.OnCalendar(now)
	want := "2026-09-12 14:00:00"
	if got != want {
		t.Errorf("cycle: got %q, want %q", got, want)
	}
}

func TestOnCalendarManual(t *testing.T) {
	s := Schedule{Kind: KindManual}
	got := s.OnCalendar(time.Now())
	if got != "" {
		t.Errorf("manual: got %q, want empty", got)
	}
}

func TestParseDDMM(t *testing.T) {
	dom, month, err := ParseDDMM("15/09")
	if err != nil {
		t.Fatalf("ParseDDMM: %v", err)
	}
	if dom != 15 || month != 9 {
		t.Errorf("ParseDDMM: got %d/%d, want 15/9", dom, month)
	}
}

func TestParseHHMM(t *testing.T) {
	hour, minute, err := ParseHHMM("21:30")
	if err != nil {
		t.Fatalf("ParseHHMM: %v", err)
	}
	if hour != 21 || minute != 30 {
		t.Errorf("ParseHHMM: got %d:%d, want 21:30", hour, minute)
	}
}

func TestParseCycle(t *testing.T) {
	cycle, err := ParseCycle("2 4 5")
	if err != nil {
		t.Fatalf("ParseCycle: %v", err)
	}
	if len(cycle) != 3 || cycle[0] != 2 || cycle[1] != 4 || cycle[2] != 5 {
		t.Errorf("ParseCycle: got %v, want [2 4 5]", cycle)
	}
}

func TestValidateDDMM(t *testing.T) {
	if err := ValidateDDMM("15/09"); err != nil {
		t.Errorf("valid DD/MM: %v", err)
	}
	if err := ValidateDDMM("32/13"); err == nil {
		t.Errorf("invalid DD/MM: should error")
	}
}

func TestValidateHHMM(t *testing.T) {
	if err := ValidateHHMM("21:30"); err != nil {
		t.Errorf("valid HH:MM: %v", err)
	}
	if err := ValidateHHMM("25:70"); err == nil {
		t.Errorf("invalid HH:MM: should error")
	}
}

func TestValidateDayOfMonth(t *testing.T) {
	if err := ValidateDayOfMonth("15"); err != nil {
		t.Errorf("valid day: %v", err)
	}
	if err := ValidateDayOfMonth("32"); err == nil {
		t.Errorf("invalid day: should error")
	}
}

func TestValidateCycle(t *testing.T) {
	if err := ValidateCycle("2 4 5"); err != nil {
		t.Errorf("valid cycle: %v", err)
	}
	if err := ValidateCycle("0 1"); err == nil {
		t.Errorf("invalid cycle: should error")
	}
}

func TestOneshotCleanupTail(t *testing.T) {
	tail := OneshotCleanupTail("/tmp/timer", "/tmp/service")
	if !strings.Contains(tail, "systemctl --user disable --now") {
		t.Error("oneshot tail missing systemctl disable")
	}
	if !strings.Contains(tail, "rm -f") {
		t.Error("oneshot tail missing rm")
	}
}

func TestCycleRescheduleTail(t *testing.T) {
	tail := CycleRescheduleTail("/tmp/recur", "/tmp/timer", "test.timer")
	if !strings.Contains(tail, "/tmp/recur") {
		t.Error("cycle tail missing recur path")
	}
	if !strings.Contains(tail, "/tmp/timer") {
		t.Error("cycle tail missing timer path")
	}
	if !strings.Contains(tail, "test.timer") {
		t.Error("cycle tail missing timer name")
	}
}

func TestWeekdayOrder(t *testing.T) {
	days := []time.Weekday{time.Friday, time.Monday, time.Wednesday}
	sorted := WeekdayOrder(days)
	if sorted[0] != time.Monday || sorted[1] != time.Wednesday || sorted[2] != time.Friday {
		t.Errorf("WeekdayOrder: got %v, want Mon,Wed,Fri", sorted)
	}
}

func TestScheduleString(t *testing.T) {
	tests := []struct {
		s    Schedule
		want string
	}{
		{Schedule{Kind: KindDaily, Hour: 21, Minute: 0}, "diário 21:00"},
		{Schedule{Kind: KindManual}, "manual"},
		{Schedule{Kind: KindOneshot, Hour: 14, Minute: 30, DOM: 15, Month: 9}, "one-shot 15/09 14:30"},
		{Schedule{Kind: KindMonthly, Hour: 8, Minute: 0, DayOfMonth: 1}, "mensal dia 1 08:00"},
		{Schedule{Kind: KindCycle, Cycle: []int{2, 4, 5}}, "ciclo 2 4 5"},
	}
	for _, tt := range tests {
		got := tt.s.String()
		if got != tt.want {
			t.Errorf("String(%v): got %q, want %q", tt.s.Kind, got, tt.want)
		}
	}
}
