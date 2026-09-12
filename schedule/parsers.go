package schedule

import (
	"fmt"
	"strconv"
	"strings"
)

// ValidateDDMM validates a DD/MM string (day 1-31, month 1-12).
func ValidateDDMM(s string) error {
	dom, month, err := ParseDDMM(s)
	if err != nil || dom < 1 || dom > 31 || month < 1 || month > 12 {
		return fmt.Errorf("expected format DD/MM")
	}
	return nil
}

// ValidateHHMM validates a HH:MM string (hour 0-23, minute 0-59).
func ValidateHHMM(s string) error {
	hour, minute, err := ParseHHMM(s)
	if err != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return fmt.Errorf("expected format HH:MM")
	}
	return nil
}

// ValidateDayOfMonth validates a day-of-month string (1-31).
func ValidateDayOfMonth(s string) error {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n < 1 || n > 31 {
		return fmt.Errorf("expected a number 1-31")
	}
	return nil
}

// ValidateCycle validates a day-interval cycle string (e.g. "2 4 5").
func ValidateCycle(s string) error {
	_, err := ParseCycle(s)
	return err
}

// ParseDDMM parses a DD/MM string into day and month.
func ParseDDMM(s string) (dom, month int, err error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format")
	}
	dom, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	month, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	return dom, month, err
}

// ParseHHMM parses a HH:MM string into hour and minute.
func ParseHHMM(s string) (hour, minute int, err error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format")
	}
	hour, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	minute, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	return hour, minute, err
}

// ParseCycle parses a custom cycle string like "2 4 5" into day-interval sequence.
func ParseCycle(s string) ([]int, error) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil, fmt.Errorf(`expected one or more day counts, e.g. "2 4 5"`)
	}
	cycle := make([]int, len(fields))
	for i, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil || n < 1 {
			return nil, fmt.Errorf(`expected positive whole numbers, e.g. "2 4 5"`)
		}
		cycle[i] = n
	}
	return cycle, nil
}
