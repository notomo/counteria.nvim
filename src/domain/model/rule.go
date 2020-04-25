package model

import (
	"database/sql/driver"
	"fmt"
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
		return "TODO"
	case TaskRuleTypeInWeekdays:
		return "TODO"
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

// Weekdays :
type Weekdays []time.Weekday

// AllWeekdays :
func AllWeekdays() Weekdays {
	return Weekdays{
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

// Days :
type Days []Day

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
// TODO
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
	case TaskRuleTypeInWeekdays:
	case TaskRuleTypeNone:
		return nil
	}
	panic("unreachable: invalid rule type: " + typ)
}

// LastTime :
// TODO
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
	case TaskRuleTypeInWeekdays:
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
		if len(rule.MonthDays()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty month days")
		}
		return nil
	case TaskRuleTypeInWeekdays:
		if len(rule.Weekdays()) == 0 {
			return NewErrValidation(ErrValidationRule, "empty weekdays")
		}
		return nil
	case TaskRuleTypeNone:
		if len(rule.Periods()) > 0 || len(rule.DateTimes()) > 0 || len(rule.Dates()) > 0 || len(rule.MonthDays()) > 0 || len(rule.Weekdays()) > 0 {
			return NewErrValidation(ErrValidationRule, "should be empty")
		}
		return nil
	}
	return NewErrValidation(ErrValidationRule, "invalid type: "+typ.String())
}
