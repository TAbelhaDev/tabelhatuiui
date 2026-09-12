package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
)

// Groups returns the huh groups for the schedule form (recurrence select +
// conditional fields). Each app builds its own huh.NewForm with its local
// groups + schedule.Groups(s). Ported from tajobs forms.go:134-254.
func Groups(s Schedule) []*huh.Group {
	timeStr := ""
	if s.Hour > 0 || s.Minute > 0 {
		timeStr = fmt.Sprintf("%02d:%02d", s.Hour, s.Minute)
	}
	dateStr := ""
	if s.DOM > 0 && s.Month > 0 {
		dateStr = fmt.Sprintf("%02d/%02d", s.DOM, s.Month)
	}
	dayOfMonthStr := ""
	if s.DayOfMonth > 0 {
		dayOfMonthStr = fmt.Sprintf("%d", s.DayOfMonth)
	}
	cycleStr := ""
	if len(s.Cycle) > 0 {
		parts := make([]string, len(s.Cycle))
		for i, d := range s.Cycle {
			parts[i] = fmt.Sprintf("%d", d)
		}
		cycleStr = strings.Join(parts, " ")
	}

	kind := s.Kind
	var weekdays []time.Weekday
	if s.Kind == KindWeekly {
		weekdays = append([]time.Weekday(nil), s.Weekdays...)
	}

	return []*huh.Group{
		// Recurrence select
		huh.NewGroup(
			huh.NewSelect[Kind]().
				Key("type").
				Title("Recurrence").
				Options(
					huh.NewOption("One-shot", KindOneshot),
					huh.NewOption("Daily", KindDaily),
					huh.NewOption("Weekly", KindWeekly),
					huh.NewOption("Monthly", KindMonthly),
					huh.NewOption("Custom cycle", KindCycle),
					huh.NewOption("Manual (no schedule)", KindManual),
				).
				Value(&kind),
		),
		// One-shot: date DD/MM + time
		huh.NewGroup(
			huh.NewInput().
				Key("date").
				Title("Date (DD/MM)").
				Placeholder("17/07").
				Value(&dateStr).
				Validate(ValidateDDMM),
			huh.NewInput().
				Key("time").
				Title("Time (HH:MM)").
				Placeholder("14:00").
				Value(&timeStr).
				Validate(ValidateHHMM),
		).WithHideFunc(func() bool { return kind != KindOneshot }),
		// Daily: time only
		huh.NewGroup(
			huh.NewInput().
				Key("time").
				Title("Time (HH:MM)").
				Placeholder("14:00").
				Value(&timeStr).
				Validate(ValidateHHMM),
		).WithHideFunc(func() bool { return kind != KindDaily }),
		// Weekly: weekdays + time
		huh.NewGroup(
			huh.NewMultiSelect[time.Weekday]().
				Key("weekdays").
				Title("Days of the week").
				Options(
					huh.NewOption("Mon", time.Monday),
					huh.NewOption("Tue", time.Tuesday),
					huh.NewOption("Wed", time.Wednesday),
					huh.NewOption("Thu", time.Thursday),
					huh.NewOption("Fri", time.Friday),
					huh.NewOption("Sat", time.Saturday),
					huh.NewOption("Sun", time.Sunday),
				).
				Value(&weekdays).
				Validate(func([]time.Weekday) error {
					if kind != KindWeekly {
						return nil
					}
					if len(weekdays) == 0 {
						return fmt.Errorf("select at least one day")
					}
					return nil
				}),
			huh.NewInput().
				Key("time").
				Title("Time (HH:MM)").
				Placeholder("14:00").
				Value(&timeStr).
				Validate(ValidateHHMM),
		).WithHideFunc(func() bool { return kind != KindWeekly }),
		// Monthly: day of month + time
		huh.NewGroup(
			huh.NewInput().
				Key("dayOfMonth").
				Title("Day of month (1-31)").
				Placeholder("15").
				Value(&dayOfMonthStr).
				Validate(ValidateDayOfMonth),
			huh.NewInput().
				Key("time").
				Title("Time (HH:MM)").
				Placeholder("14:00").
				Value(&timeStr).
				Validate(ValidateHHMM),
		).WithHideFunc(func() bool { return kind != KindMonthly }),
		// Cycle: cycle pattern + first run date + time
		huh.NewGroup(
			huh.NewInput().
				Key("cycle").
				Title(`Day-interval cycle (e.g. "2 4 5")`).
				Placeholder("2 4 5").
				Value(&cycleStr).
				Validate(ValidateCycle),
			huh.NewInput().
				Key("date").
				Title("First run date (DD/MM)").
				Placeholder("17/07").
				Value(&dateStr).
				Validate(ValidateDDMM),
			huh.NewInput().
				Key("time").
				Title("First run time (HH:MM)").
				Placeholder("14:00").
				Value(&timeStr).
				Validate(ValidateHHMM),
		).WithHideFunc(func() bool { return kind != KindCycle }),
	}
}
