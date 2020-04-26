package model

import (
	"database/sql/driver"
	"fmt"
	"math"
	"time"
)

// TaskRule :
type TaskRule struct {
	TaskRuleData
}

func (rule *TaskRule) String() string {
	typ := rule.Type()
	switch typ {
	case TaskRuleTypePeriodic:
		period := rule.Periods()[0]
		return fmt.Sprintf("once per %d %s", period.Number(), period.Unit())
	case TaskRuleTypeByTimes:
		dt := rule.DateTimes()[0]
		return fmt.Sprintf("by %s", dt.Format("2006-01-02 15:04:05"))
	case TaskRuleTypeInDates:
		date := rule.Dates()[0]
		return fmt.Sprintf("in %s", date)
	case TaskRuleTypeInDaysEveryMonth:
		day := rule.Days()[0]
		return fmt.Sprintf("in %d evety month", day)
	case TaskRuleTypeInWeekdays:
		weekday := rule.Weekdays()[0]
		return fmt.Sprintf("in %s", weekday)
	case TaskRuleTypeNone:
		return "None"
	}
	panic("invalid rule type: " + typ)
}

// TaskRuleData :
type TaskRuleData interface {
	Type() TaskRuleType

	Weekdays() Weekdays
	Dates() Dates
	MonthDays() MonthDays
	Days() Days
	DateTimes() DateTimes
	Periods() Periods
}

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

// MonthDay : mm-dd
type MonthDay string

// MonthDays :
type MonthDays []MonthDay

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

// DateTimes :
type DateTimes []time.Time

// NextTime :
func (dt DateTimes) NextTime(at time.Time) *time.Time {
	for _, t := range dt {
		if t.After(at) {
			return &t
		}
	}
	return nil
}

// TaskRuleType :
type TaskRuleType string

var (
	// TaskRuleTypePeriodic :
	TaskRuleTypePeriodic = TaskRuleType("periodic")
	// TaskRuleTypeByTimes : oneshot deadline by time
	TaskRuleTypeByTimes = TaskRuleType("byTimes")
	// TaskRuleTypeInDaysEveryMonth :
	TaskRuleTypeInDaysEveryMonth = TaskRuleType("inDaysEveryMonth")
	// TaskRuleTypeInDates :
	TaskRuleTypeInDates = TaskRuleType("inDates")
	// TaskRuleTypeInWeekdays :
	TaskRuleTypeInWeekdays = TaskRuleType("inWeekdays")
	// TaskRuleTypeNone :
	TaskRuleTypeNone = TaskRuleType("none")
)

func (typ TaskRuleType) String() string {
	return string(typ)
}

// TaskRuleTypes :
func TaskRuleTypes() []TaskRuleType {
	return []TaskRuleType{
		TaskRuleTypePeriodic,
		TaskRuleTypeByTimes,
		TaskRuleTypeInDaysEveryMonth,
		TaskRuleTypeInDates,
		TaskRuleTypeInWeekdays,
		TaskRuleTypeNone,
	}
}

// NextTime :
func (rule *TaskRule) NextTime(startAt time.Time, lastDone *DoneTask) *time.Time {
	typ := rule.Type()
	switch typ {
	case TaskRuleTypePeriodic:
		if lastDone == nil {
			return rule.Periods().NextTime(startAt)
		}
		return rule.Periods().NextTime(lastDone.At())
	case TaskRuleTypeByTimes:
		if lastDone == nil {
			return rule.DateTimes().NextTime(startAt)
		}
		return nil
	case TaskRuleTypeInDates:
		if lastDone == nil {
			return rule.Dates().NextTime(startAt)
		}
		return nil
	case TaskRuleTypeInDaysEveryMonth:
		if lastDone == nil {
			return rule.Days().NextTime(startAt)
		}
		return rule.Days().NextTime(lastDone.At())
	case TaskRuleTypeInWeekdays:
		if lastDone == nil {
			return rule.Weekdays().NextTime(startAt)
		}
		return rule.Weekdays().NextTime(lastDone.At())
	case TaskRuleTypeNone:
		return nil
	}
	panic("unreachable: invalid rule type: " + typ)
}

// LastTime :
func (rule *TaskRule) LastTime(startAt time.Time, lastDone *DoneTask) *time.Time {
	typ := rule.Type()
	switch typ {
	case TaskRuleTypePeriodic:
		return nil
	case TaskRuleTypeByTimes:
		return rule.DateTimes().NextTime(startAt)
	case TaskRuleTypeInDates:
		return rule.Dates().NextTime(startAt)
	case TaskRuleTypeInDaysEveryMonth:
		if lastDone == nil {
			return rule.Days().NextTime(startAt)
		}
		return rule.Days().NextTime(lastDone.At())
	case TaskRuleTypeInWeekdays:
		if lastDone == nil {
			return rule.Weekdays().NextTime(startAt)
		}
		return rule.Weekdays().NextTime(lastDone.At())
	case TaskRuleTypeNone:
		return nil
	}
	panic("unreachable: invalid rule type: " + typ)
}

// Validate :
func (rule *TaskRule) Validate() error {
	typ := rule.Type()
	switch typ {
	case TaskRuleTypePeriodic:
		if len(rule.Periods()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty periods")
		}
		return nil
	case TaskRuleTypeByTimes:
		if len(rule.DateTimes()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty date times")
		}
		return nil
	case TaskRuleTypeInDates:
		if len(rule.Dates()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty dates")
		}
		return nil
	case TaskRuleTypeInDaysEveryMonth:
		if len(rule.Days()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty days")
		}
		return nil
	case TaskRuleTypeInWeekdays:
		if len(rule.Weekdays()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty weekdays")
		}
		return nil
	case TaskRuleTypeNone:
		if len(rule.Periods()) > 0 || len(rule.DateTimes()) > 0 || len(rule.Dates()) > 0 || len(rule.Days()) > 0 || len(rule.MonthDays()) > 0 || len(rule.Weekdays()) > 0 {
			return NewErrValidation(ErrValidationRule, "should be empty")
		}
		return nil
	}
	return NewErrValidation(ErrValidationRule, "invalid type: "+typ.String())
}
