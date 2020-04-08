package lib

import "time"

// Clock :
type Clock interface {
	Now() time.Time
}

// NewClock : default impl
func NewClock() Clock {
	return &clock{}
}

type clock struct{}

func (c *clock) Now() time.Time {
	return time.Now()
}
