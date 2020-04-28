package model

import (
	"database/sql/driver"
	"time"
)

// Date : yyyy-mm-dd
type Date string

const dateFormat = "2006-01-02"

// Time : date to time
func (date Date) Time() time.Time {
	t, _ := time.Parse(dateFormat, string(date))
	return t
}

// Contains :
func (date Date) Contains(at time.Time) bool {
	t := date.Time()
	y, m, d := t.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	end := time.Date(y, m, d, 23, 59, 59, 0, time.Local)
	return start.Before(at) && end.After(at)
}

// Value : FIXME: for datestore
func (date Date) Value() (driver.Value, error) {
	return driver.Value(date.Time()), nil
}

// Scan : FIXME: for datestore
func (date *Date) Scan(value interface{}) error {
	*date = Date(string(value.(time.Time).Format(dateFormat)))
	return nil
}

// Dates :
type Dates []Date

// NextTime :
func (dates Dates) NextTime(at time.Time) *time.Time {
	for _, d := range dates {
		t := d.Time()
		if t.After(at) {
			return &t
		}
	}
	return nil
}

// Contains :
func (dates Dates) Contains(at time.Time) bool {
	for _, d := range dates {
		if d.Contains(at) {
			return true
		}
	}
	return false
}
