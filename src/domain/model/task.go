package model

import (
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
	Rule() TaskRule
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

// Deadline :
func (task *Task) Deadline() Deadline {
	return Deadline{
		Rule:     task.Rule(),
		StartAt:  task.StartAt(),
		LastDone: task.LastDone(),
	}
}

// Deadline :
type Deadline struct {
	Rule     TaskRule
	StartAt  time.Time
	LastDone *DoneTask
}

// Next :
func (deadline Deadline) Next() *time.Time {
	return deadline.Rule.NextTime(deadline.StartAt, deadline.LastDone)
}

// Latest :
func (deadline Deadline) Latest() time.Time {
	next := deadline.Next()
	if next != nil {
		return *next
	}
	return *deadline.Rule.LastTime(deadline.StartAt, deadline.LastDone)
}

// RemainingTime : how much time until task deadline
func (deadline Deadline) RemainingTime(now time.Time) RemainingTime {
	duration := deadline.Latest().Sub(now)

	h := int(duration.Hours())
	days := h / 24
	hours := h % 24
	minutes := int(duration.Minutes()) % 60

	return RemainingTime{
		Days:     days,
		Hours:    hours,
		Minutes:  minutes,
		duration: duration,
	}
}

// RemainingTime :
type RemainingTime struct {
	Days    int
	Hours   int
	Minutes int

	duration time.Duration
}

// Exists :
func (t RemainingTime) Exists() bool {
	return t.duration > 0
}

// DoneTask :
type DoneTask struct {
	DoneTaskData
}

// DoneTaskData :
type DoneTaskData interface {
	At() time.Time
}
