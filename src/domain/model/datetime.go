package model

import "time"

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
