package model

import "time"

// Period :
type Period struct {
	PeriodData
}

// PeriodData :
type PeriodData interface {
	Number() int
	Unit() PeriodUnit
}

// FromTime : return from + period
func (period Period) FromTime(from time.Time) time.Time {
	year, month, day := period.Unit().numbers()
	number := period.Number()
	return from.AddDate(year*number, month*number, day*number)
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

func (unit PeriodUnit) numbers() (year int, month int, day int) {
	switch unit {
	case PeriodUnitYear:
		year = 1
		return
	case PeriodUnitMonth:
		month = 1
		return
	case PeriodUnitWeek:
		day = 7
		return
	case PeriodUnitDay:
		day = 1
		return
	}
	panic("unreachable: invalid period unit: " + unit)
}

func (unit PeriodUnit) String() string {
	return string(unit)
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

// Periods :
type Periods []Period

// NextTime :
func (periods Periods) NextTime(at time.Time) *time.Time {
	for _, p := range periods {
		t := p.FromTime(at)
		return &t
	}
	return nil
}
