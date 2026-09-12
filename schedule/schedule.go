package schedule

import (
	"fmt"
	"strings"
	"time"
)

// Kind is the type of recurrence for a schedule.
type Kind int

const (
	KindOneshot Kind = iota
	KindDaily
	KindWeekly
	KindMonthly
	KindCycle
	KindManual
)

// Schedule is the complete definition of a schedule.
type Schedule struct {
	Kind       Kind
	Minute     int
	Hour       int
	DOM        int    // oneshot/cycle: first run day-of-month (1-31)
	Month      int    // oneshot/cycle: first run month (1-12)
	Weekdays   []time.Weekday // weekly
	DayOfMonth int            // monthly (1-31)
	Cycle      []int          // cycle: day-interval sequence
}

// String returns a human-readable summary for the metadata panel.
func (s Schedule) String() string {
	switch s.Kind {
	case KindOneshot:
		return fmt.Sprintf("one-shot %02d/%02d %02d:%02d", s.DOM, s.Month, s.Hour, s.Minute)
	case KindDaily:
		return fmt.Sprintf("diário %02d:%02d", s.Hour, s.Minute)
	case KindWeekly:
		sorted := WeekdayOrder(s.Weekdays)
		abbrs := make([]string, len(sorted))
		for i, d := range sorted {
			abbrs[i] = WeekdayAbbr[d]
		}
		return fmt.Sprintf("semanal %s %02d:%02d", strings.Join(abbrs, ","), s.Hour, s.Minute)
	case KindMonthly:
		return fmt.Sprintf("mensal dia %d %02d:%02d", s.DayOfMonth, s.Hour, s.Minute)
	case KindCycle:
		parts := make([]string, len(s.Cycle))
		for i, d := range s.Cycle {
			parts[i] = fmt.Sprintf("%d", d)
		}
		return fmt.Sprintf("ciclo %s", strings.Join(parts, " "))
	case KindManual:
		return "manual"
	default:
		return "desconhecido"
	}
}
