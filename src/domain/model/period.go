package model

import "time"

// Period :
type Period interface {
	Number() int
	Unit() PeriodUnit
	FromTime(time.Time) time.Time
}

// PeriodUnit :
type PeriodUnit string

var (
	// PeriodUnitYear :
	PeriodUnitYear = PeriodUnit("year")
	// PeriodUnitMonth :
	PeriodUnitMonth = PeriodUnit("month")
	// PeriodUnitWeek :
	PeriodUnitWeek = PeriodUnit("week")
	// PeriodUnitDay :
	PeriodUnitDay = PeriodUnit("day")
)

// Numbers : year, month, day
func (unit PeriodUnit) Numbers() (int, int, int) {
	switch unit {
	case PeriodUnitYear:
		return 1, 0, 0
	case PeriodUnitMonth:
		return 0, 1, 0
	case PeriodUnitWeek:
		return 0, 0, 7
	case PeriodUnitDay:
		return 0, 0, 1
	}
	panic("unreachable: invalid period unit: " + unit)
}

// PeriodUnits :
func PeriodUnits() []PeriodUnit {
	return []PeriodUnit{
		PeriodUnitYear,
		PeriodUnitMonth,
		PeriodUnitWeek,
		PeriodUnitDay,
	}
}
