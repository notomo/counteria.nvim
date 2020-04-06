package model

import "time"

// Task :
type Task interface {
	ID() int
	Name() string
	Period() Period
	DoneAt() *time.Time
	LimitAt() time.Time
	PastDeadline(time.Time) bool
}
