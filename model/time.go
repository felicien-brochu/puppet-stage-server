package model

import "math"

// Duration number of nanoseconds (max: 290 years)
type Duration int64

// Time number of nanoseconds (max: 290 years)
type Time int64

// After returns returns true if t is after t1
func (t Time) After(t1 Time) bool {
	return t >= t1
}

// Before returns returns true if t is before t1
func (t Time) Before(t1 Time) bool {
	return t <= t1
}

// Since returns the Duration between t1 and t
func (t Time) Since(t1 Time) Duration {
	return Duration(t - t1)
}

// AbsDiff returns the absolute difference between two times
func (t Time) AbsDiff(t1 Time) Duration {
	return Duration(math.Abs(float64(t - t1)))
}

// Common durations
const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)
