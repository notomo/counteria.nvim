package sqliteimpl

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/pkg/errors"
)

// TaskRuleLineRepository :
type TaskRuleLineRepository struct {
	Db *gorp.DbMap
}

// Create :
func (repo *TaskRuleLineRepository) Create(transaction repository.Transaction, taskID int, rule *model.TaskRule) error {
	trans := transaction.(*gorp.Transaction)

	lines := repo.toRuleLines(taskID, rule)
	ls := make([]interface{}, len(lines))
	for i, l := range lines {
		l := l
		ls[i] = &l
	}

	if err := trans.Insert(ls...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete :
func (repo *TaskRuleLineRepository) Delete(transaction repository.Transaction, taskID int) error {
	trans := transaction.(*gorp.Transaction)

	lines, err := repo.List(taskID)
	if err != nil {
		return errors.WithStack(err)
	}

	ls := make([]interface{}, len(lines))
	for i, l := range lines {
		l := l
		ls[i] = &l
	}

	if _, err := trans.Delete(ls...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// List :
func (repo *TaskRuleLineRepository) List(taskIDs ...int) ([]TaskRuleLine, error) {
	lines := []TaskRuleLine{}
	if len(taskIDs) == 0 {
		return lines, nil
	}

	if _, err := repo.Db.Select(&lines, `
	SELECT *
	FROM task_rule_lines
	WHERE task_id IN (:ids)
	`, map[string]interface{}{"ids": taskIDs}); err != nil {
		return nil, errors.WithStack(err)
	}
	return lines, nil
}

// Update : delete and insert
func (repo *TaskRuleLineRepository) Update(transaction repository.Transaction, taskID int, rule *model.TaskRule) error {
	trans := transaction.(*gorp.Transaction)

	if err := repo.Delete(trans, taskID); err != nil {
		return errors.WithStack(err)
	}

	if err := repo.Create(trans, taskID, rule); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *TaskRuleLineRepository) toRuleLines(taskID int, rule *model.TaskRule) []TaskRuleLine {
	lines := []TaskRuleLine{}
	switch typ := rule.TaskRuleData.Type(); typ {
	case model.TaskRuleTypePeriodic:
		for _, p := range rule.Periods() {
			number := p.Number()
			unit := p.Unit()
			line := TaskRuleLine{
				TaskID: taskID,
				TaskPeriod: TaskPeriod{
					PeriodNumber: &number,
					PeriodUnit:   &unit,
				},
			}
			lines = append(lines, line)
		}
	case model.TaskRuleTypeByTimes:
	case model.TaskRuleTypeInDaysEveryMonth:
	case model.TaskRuleTypeInDates:
	case model.TaskRuleTypeInWeekdays:
	default:
		panic("unreachable: invalid rule type: " + typ)
	}
	return lines
}

var _ model.TaskRuleData = &TaskRule{}

// TaskRule :
type TaskRule struct {
	RuleType      model.TaskRuleType
	RuleWeekdays  model.Weekdays
	RuleDays      model.Days
	RuleMonthDays model.MonthDays
	RuleTimes     model.Times
	RuleDates     model.Dates
	RulePeriods   model.Periods
}

func (rule *TaskRule) addLine(line TaskRuleLine) {
	typ := rule.RuleType
	switch typ {
	case model.TaskRuleTypePeriodic:
		rule.RulePeriods = append(rule.RulePeriods, model.Period{
			PeriodData: &TaskPeriod{
				PeriodNumber: line.PeriodNumber,
				PeriodUnit:   line.PeriodUnit,
			},
		})
	case model.TaskRuleTypeByTimes:
	case model.TaskRuleTypeInDates:
	case model.TaskRuleTypeInDaysEveryMonth:
	case model.TaskRuleTypeInWeekdays:
	default:
		panic("invalid rule type: " + typ)
	}
}

// Type :
func (rule *TaskRule) Type() model.TaskRuleType {
	return rule.RuleType
}

// Weekdays :
func (rule *TaskRule) Weekdays() model.Weekdays {
	return rule.RuleWeekdays
}

// Days :
func (rule *TaskRule) Days() model.Days {
	return rule.RuleDays
}

// MonthDays :
func (rule *TaskRule) MonthDays() model.MonthDays {
	return rule.RuleMonthDays
}

// Dates :
func (rule *TaskRule) Dates() model.Dates {
	return rule.RuleDates
}

// DateTimes :
func (rule *TaskRule) DateTimes() model.Times {
	return rule.RuleTimes
}

// Periods :
func (rule *TaskRule) Periods() model.Periods {
	return rule.RulePeriods
}

// TaskRuleLine :
type TaskRuleLine struct {
	ID       int             `db:"id, primarykey, autoincrement"`
	TaskID   int             `db:"task_id, notnull" foreign:"tasks(id)"`
	Weekday  *time.Weekday   `db:"weekday" check:"weekday"`
	Day      *model.Day      `db:"day" check:"day"`
	MonthDay *model.MonthDay `db:"month_day"`
	Time     *time.Time      `db:"time"`
	Date     *model.Date     `db:"date"`
	TaskPeriod
}

var _ model.PeriodData = &TaskPeriod{}

// TaskPeriod :
type TaskPeriod struct {
	PeriodNumber *int              `db:"period_number" check:"natural"`
	PeriodUnit   *model.PeriodUnit `db:"period_unit" check:"periodUnit"`
}

// Number :
func (period TaskPeriod) Number() int {
	return *period.PeriodNumber
}

// Unit :
func (period TaskPeriod) Unit() model.PeriodUnit {
	return *period.PeriodUnit
}
