package model

import (
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
		day := rule.Days()[0]
		return fmt.Sprintf("in %d every month", day)
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
