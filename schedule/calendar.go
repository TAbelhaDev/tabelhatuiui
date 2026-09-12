package schedule

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// WeekdayAbbr maps time.Weekday to systemd day abbreviations.
var WeekdayAbbr = map[time.Weekday]string{
	time.Monday:    "Mon",
	time.Tuesday:   "Tue",
	time.Wednesday: "Wed",
	time.Thursday:  "Thu",
	time.Friday:    "Fri",
	time.Saturday:  "Sat",
	time.Sunday:    "Sun",
}

// isoWeekday returns weekday number with Monday=1..Sunday=7.
func isoWeekday(d time.Weekday) int {
	if d == time.Sunday {
		return 7
	}
	return int(d)
}

// WeekdayOrder sorts weekdays into Mon..Sun display order.
func WeekdayOrder(weekdays []time.Weekday) []time.Weekday {
	sorted := append([]time.Weekday(nil), weekdays...)
	sort.Slice(sorted, func(i, k int) bool {
		return isoWeekday(sorted[i]) < isoWeekday(sorted[k])
	})
	return sorted
}

// OnCalendar generates the systemd OnCalendar expression for this schedule.
func (s Schedule) OnCalendar(now time.Time) string {
	switch s.Kind {
	case KindDaily:
		return fmt.Sprintf("*-*-* %02d:%02d:00", s.Hour, s.Minute)
	case KindWeekly:
		sorted := WeekdayOrder(s.Weekdays)
		abbrs := make([]string, len(sorted))
		for i, d := range sorted {
			abbrs[i] = WeekdayAbbr[d]
		}
		return fmt.Sprintf("%s *-*-* %02d:%02d:00", strings.Join(abbrs, ","), s.Hour, s.Minute)
	case KindMonthly:
		return fmt.Sprintf("*-*-%02d %02d:%02d:00", s.DayOfMonth, s.Hour, s.Minute)
	case KindOneshot, KindCycle:
		return computeAbsolute(s.Minute, s.Hour, s.DOM, s.Month, now)
	case KindManual:
		return ""
	default:
		return ""
	}
}

// computeAbsolute builds a systemd OnCalendar= timestamp for one-shot/cycle.
// Year rolls over to next year if the given month/day already passed.
func computeAbsolute(minute, hour, dom, month int, now time.Time) string {
	year := now.Year()
	todayMonth, todayDay := int(now.Month()), now.Day()
	if month < todayMonth || (month == todayMonth && dom < todayDay) {
		year++
	}
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", year, month, dom, hour, minute)
}
