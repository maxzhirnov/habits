package utils

import "time"

// DateOnly returns date with zeroed hours, minute, seconds and nanoseconds
func DateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
