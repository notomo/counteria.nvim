package model

import (
	"database/sql/driver"
	"time"
)

// TaskRule :
type TaskRule struct {
	TaskRuleData
}

// TaskRuleData :
type TaskRuleData interface {
	Type() TaskRuleType

	Weekdays() Weekdays
	Dates() Dates
	MonthDays() MonthDays
	Days() Days
	DateTimes() Times
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

// MonthDay : mm-dd
type MonthDay string

// MonthDays :
type MonthDays []MonthDay

// Day : dd
type Day int

// Days :
type Days []Day

// Times :
type Times []time.Time

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
)

// TaskRuleTypes :
func TaskRuleTypes() []TaskRuleType {
	return []TaskRuleType{
		TaskRuleTypePeriodic,
		TaskRuleTypeByTimes,
		TaskRuleTypeInDaysEveryMonth,
		TaskRuleTypeInDates,
		TaskRuleTypeInWeekdays,
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
	}
	panic("unreachable: invalid rule type: " + typ)
}

// LastTime :
// TODO
func (rule *TaskRule) LastTime(startAt time.Time, lastDone *DoneTask) *time.Time {
	typ := rule.Type()
	switch typ {
	case TaskRuleTypePeriodic:
		if lastDone == nil {
			return rule.Periods().NextTime(startAt)
		}
		return rule.Periods().NextTime(lastDone.At())
	case TaskRuleTypeByTimes:
	}
	panic("unreachable: invalid rule type: " + typ)
}
