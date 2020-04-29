package model

import "time"

// Day : dd
type Day int

// Contains :
func (day Day) Contains(at time.Time) bool {
	d := at.Day()
	return int(day) == d
}

// NextTime :
func (day Day) NextTime(at time.Time) time.Time {
	y, m, d := at.Date()
	targetDay := int(day)
	if targetDay < d {
		m = m + 1
	}
	t := time.Date(y, m, targetDay, 23, 59, 59, 999999999, time.Local)
	if t.Month() == m {
		return t
	}
	return time.Date(t.Year(), t.Month(), 1, 23, 59, 59, 999999999, time.Local).AddDate(0, 0, -1)
}

// Days :
type Days []Day

// Contains :
func (days Days) Contains(at time.Time) bool {
	for _, d := range days {
		if d.Contains(at) {
			return true
		}
	}
	return false
}

// NextTime :
func (days Days) NextTime(at time.Time) *time.Time {
	for _, d := range days {
		t := d.NextTime(at)
		return &t
	}
	return nil
}