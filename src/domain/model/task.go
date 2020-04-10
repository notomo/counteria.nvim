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
	Period() Period
	StartAt() time.Time
	LastDone() *DoneTask
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

// LimitAt :
func (task *Task) LimitAt() time.Time {
	lastDone := task.LastDone()
	if lastDone == nil {
		return task.Period().FromTime(task.StartAt())
	}
	return task.Period().FromTime(lastDone.At())
}

// PastDeadline :
func (task *Task) PastDeadline(now time.Time) bool {
	return now.After(task.LimitAt())
}

// RemainingTime : how much time until task deadline
func (task *Task) RemainingTime(now time.Time) RemainingTime {
	duration := task.LimitAt().Sub(now)

	h := int(duration.Hours())
	days := h / 24
	hours := h % 24
	minutes := int(duration.Minutes()) % 60

	return RemainingTime{
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
	}
}

// RemainingTime :
type RemainingTime struct {
	Days    int
	Hours   int
	Minutes int
}

// DoneTask :
type DoneTask struct {
	DoneTaskData
}

// DoneTaskData :
type DoneTaskData interface {
	At() time.Time
}
