package model

// Period :
type Period interface {
	Number() int
	Unit() PeriodUnit
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
