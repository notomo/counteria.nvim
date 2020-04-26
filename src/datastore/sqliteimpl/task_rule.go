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
func (repo *TaskRuleLineRepository) Create(transaction repository.Transaction, task *Task) error {
	trans := transaction.(*gorp.Transaction)

	lines := task.ruleLines()
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

// Bind :
func (repo *TaskRuleLineRepository) Bind(tasks ...*Task) error {
	ids := make([]int, len(tasks))
	taskMap := make(map[int]*Task)
	for i, task := range tasks {
		ids[i] = task.TaskID
		taskMap[task.TaskID] = task
		task.TaskRule = NewTaskRule(task.TaskRuleType)
	}

	lines, err := repo.List(ids...)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, line := range lines {
		task := taskMap[line.TaskID]
		task.TaskRule.add(line)
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
func (repo *TaskRuleLineRepository) Update(transaction repository.Transaction, task *Task) error {
	trans := transaction.(*gorp.Transaction)

	if err := repo.Delete(trans, task.ID()); err != nil {
		return errors.WithStack(err)
	}

	if err := repo.Create(trans, task); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var _ model.TaskRuleData = &TaskRule{}

// NewTaskRule :
func NewTaskRule(typ model.TaskRuleType, opts ...func(*TaskRule)) *TaskRule {
	rule := &TaskRule{
		RuleType:      typ,
		RuleWeekdays:  model.Weekdays{},
		RuleDays:      model.Days{},
		RuleMonthDays: model.MonthDays{},
		RuleDateTimes: model.DateTimes{},
		RuleDates:     model.Dates{},
		RulePeriods:   model.Periods{},
	}
	for _, opt := range opts {
		opt(rule)
	}
	return rule
}

// WithPeriod :
func WithPeriod(number int, unit model.PeriodUnit) func(*TaskRule) {
	return func(ob *TaskRule) {
		period := model.Period{
			PeriodData: TaskPeriod{
				PeriodNumber: &number,
				PeriodUnit:   &unit,
			},
		}
		ob.RulePeriods = append(ob.RulePeriods, period)
	}
}

// TaskRule :
type TaskRule struct {
	RuleType      model.TaskRuleType
	RuleWeekdays  model.Weekdays
	RuleDays      model.Days
	RuleMonthDays model.MonthDays
	RuleDateTimes model.DateTimes
	RuleDates     model.Dates
	RulePeriods   model.Periods
}

func (rule *TaskRule) add(line TaskRuleLine) {
	typ := rule.RuleType
	switch typ {
	case model.TaskRuleTypePeriodic:
		rule.RulePeriods = append(rule.RulePeriods, model.Period{
			PeriodData: &TaskPeriod{
				PeriodNumber: line.PeriodNumber,
				PeriodUnit:   line.PeriodUnit,
			},
		})
		return
	case model.TaskRuleTypeByTimes:
		rule.RuleDateTimes = append(rule.RuleDateTimes, *line.DateTime)
		return
	case model.TaskRuleTypeInDates:
		rule.RuleDates = append(rule.RuleDates, *line.Date)
		return
	case model.TaskRuleTypeInDaysEveryMonth:
		rule.RuleDays = append(rule.RuleDays, *line.Day)
		return
	case model.TaskRuleTypeInWeekdays:
		rule.RuleWeekdays = append(rule.RuleWeekdays, *line.Weekday)
		return
	case model.TaskRuleTypeNone:
		return
	}
	panic("invalid rule type: " + typ)
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
func (rule *TaskRule) DateTimes() model.DateTimes {
	return rule.RuleDateTimes
}

// Periods :
func (rule *TaskRule) Periods() model.Periods {
	return rule.RulePeriods
}

// TaskRuleLine :
type TaskRuleLine struct {
	ID       int             `db:"id, primarykey, autoincrement"`
	TaskID   int             `db:"task_id, notnull" foreign:"tasks(id)"`
	Weekday  *model.Weekday  `db:"weekday" check:"weekday"`
	Day      *model.Day      `db:"day" check:"day"`
	MonthDay *model.MonthDay `db:"month_day"`
	DateTime *time.Time      `db:"date_time"`
	Date     *model.Date     `db:"rule_date"` // avoid using `date`
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

func (task *Task) ruleLines() []TaskRuleLine {
	lines := []TaskRuleLine{}
	typ := task.TaskRuleType
	switch typ {
	case model.TaskRuleTypePeriodic:
		for _, p := range task.Rule().Periods() {
			number := p.Number()
			unit := p.Unit()
			lines = append(lines, TaskRuleLine{
				TaskID: task.ID(),
				TaskPeriod: TaskPeriod{
					PeriodNumber: &number,
					PeriodUnit:   &unit,
				},
			})
		}
		return lines
	case model.TaskRuleTypeByTimes:
		for _, t := range task.Rule().DateTimes() {
			t := t
			lines = append(lines, TaskRuleLine{
				TaskID:   task.ID(),
				DateTime: &t,
			})
		}
		return lines
	case model.TaskRuleTypeInDaysEveryMonth:
		for _, day := range task.Rule().Days() {
			day := day
			lines = append(lines, TaskRuleLine{
				TaskID: task.ID(),
				Day:    &day,
			})
		}
		return lines
	case model.TaskRuleTypeInDates:
		for _, date := range task.Rule().Dates() {
			date := date
			lines = append(lines, TaskRuleLine{
				TaskID: task.ID(),
				Date:   &date,
			})
		}
		return lines
	case model.TaskRuleTypeInWeekdays:
		for _, weekday := range task.Rule().Weekdays() {
			weekday := weekday
			lines = append(lines, TaskRuleLine{
				TaskID:  task.ID(),
				Weekday: &weekday,
			})
		}
		return lines
	case model.TaskRuleTypeNone:
		return lines
	}
	panic("unreachable: invalid rule type: " + typ)
}
