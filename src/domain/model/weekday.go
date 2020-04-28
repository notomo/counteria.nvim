package model

import (
	"math"
	"time"
)

// Weekday :
type Weekday time.Weekday

// NextTime :
func (weekday Weekday) NextTime(at time.Time) time.Time {
	w := at.Weekday()
	diff := int(math.Abs(float64(w - time.Weekday(weekday))))
	y, m, d := at.Date()
	return time.Date(y, m, d+diff, 23, 59, 59, 999999999, time.Local)
}

// Contains :
func (weekday Weekday) Contains(at time.Time) bool {
	w := at.Weekday()
	return w == time.Weekday(weekday)
}

func (weekday Weekday) String() string {
	return time.Weekday(weekday).String()
}

// Weekdays :
type Weekdays []Weekday

// NextTime :
func (weekdays Weekdays) NextTime(at time.Time) *time.Time {
	for _, d := range weekdays {
		t := d.NextTime(at)
		return &t
	}
	return nil
}

// Contains :
func (weekdays Weekdays) Contains(at time.Time) bool {
	for _, d := range weekdays {
		if d.Contains(at) {
			return true
		}
	}
	return false
}

// AllWeekdays :
func AllWeekdays() []time.Weekday {
	return []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}
}
