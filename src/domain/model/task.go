package model

import "time"

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

// DoneAt :
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

// DoneTask :
type DoneTask struct {
	DoneTaskData
}

// DoneTaskData :
type DoneTaskData interface {
	At() time.Time
}
