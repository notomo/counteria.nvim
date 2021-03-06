package model

import (
	"math"
	"time"
)

// Task :
type Task struct {
	TaskData
}

// TaskData :
type TaskData interface {
	ID() int
	Name() string
	StartAt() time.Time
	LastDone() *DoneTask
	Rule() *TaskRule
}

// Validate :
func (task *Task) Validate() error {
	return task.Rule().Validate()
}

// DoneAt : the time the task was done
func (task *Task) DoneAt() *time.Time {
	lastDone := task.LastDone()
	if lastDone == nil {
		return nil
	}
	at := lastDone.At()
	return &at
}

// Done :
func (task *Task) Done(now time.Time) bool {
	typ := task.Rule().Type()
	switch typ {
	case TaskRuleTypePeriodic:
		return false
	case TaskRuleTypeByTimes:
		return task.LastDone() != nil
	case TaskRuleTypeInDaysEveryMonth:
		return task.LastDone() != nil && task.Rule().Days().Contains(task.LastDone().At(), now)
	case TaskRuleTypeInDates:
		return task.LastDone() != nil
	case TaskRuleTypeInWeekdays:
		return task.LastDone() != nil && task.Rule().Weekdays().Contains(task.LastDone().At())
	case TaskRuleTypeNone:
		return task.LastDone() != nil
	}
	panic("unreachable: invalid rule type: " + typ)
}

// IsActive :
func (task *Task) IsActive(now time.Time) bool {
	rule := task.Rule()
	typ := rule.Type()
	switch typ {
	case TaskRuleTypePeriodic:
		return true
	case TaskRuleTypeByTimes:
		return true
	case TaskRuleTypeInDaysEveryMonth:
		return rule.Days().Contains(now, now)
	case TaskRuleTypeInDates:
		return rule.Dates().Contains(now)
	case TaskRuleTypeInWeekdays:
		return rule.Weekdays().Contains(now)
	case TaskRuleTypeNone:
		return true
	}
	panic("unreachable: invalid rule type: " + typ)
}

// Deadline :
func (task *Task) Deadline(now time.Time) Deadline {
	return Deadline{
		Rule:     task.Rule(),
		StartAt:  task.StartAt(),
		LastDone: task.LastDone(),
		Done:     task.Done(now),
		Now:      now,
	}
}

// Deadline :
type Deadline struct {
	Rule     *TaskRule
	StartAt  time.Time
	LastDone *DoneTask
	Done     bool
	Now      time.Time
}

// Next :
func (deadline Deadline) Next() *time.Time {
	return deadline.Rule.NextTime(deadline.StartAt, deadline.LastDone)
}

// Latest :
func (deadline Deadline) Latest() *time.Time {
	next := deadline.Next()
	if next != nil {
		return next
	}
	return deadline.Rule.LastTime(deadline.StartAt, deadline.LastDone)
}

// RemainingTime : how much time until task deadline
func (deadline Deadline) RemainingTime() RemainingTime {
	latest := deadline.Latest()
	if latest == nil {
		return RemainingTime{Done: deadline.Done}
	}
	duration := latest.Sub(deadline.Now)

	h := int(math.Abs(duration.Hours()))
	days := h / 24
	hours := h % 24
	minutes := int(math.Abs(duration.Minutes())) % 60

	return RemainingTime{
		Days:     days,
		Hours:    hours,
		Minutes:  minutes,
		Done:     deadline.Done,
		duration: duration,
	}
}

// RemainingTime :
type RemainingTime struct {
	Days    int
	Hours   int
	Minutes int

	Done     bool
	duration time.Duration
}

// Exists :
func (t RemainingTime) Exists() bool {
	return t.duration >= 0
}

// DoneTask :
type DoneTask struct {
	DoneTaskData
}

// DoneTaskData :
type DoneTaskData interface {
	At() time.Time
}
