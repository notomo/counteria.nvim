package sqliteimpl

import (
	"database/sql"
	"encoding/json"
	"io"
	"sort"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/pkg/errors"
)

// TaskRepository : impl
type TaskRepository struct {
	Db *gorp.DbMap
}

var _ repository.TaskRepository = &TaskRepository{}

// TaskSummary :
type TaskSummary struct {
	TaskID       int                `json:"id" db:"id"`
	TaskName     string             `json:"name" db:"name"`
	StartAt      time.Time          `json:"startAt" db:"start_at"`
	TaskRuleType model.TaskRuleType `json:"ruleType" db:"rule_type" check:"taskRuleType"`

	LastDoneID *int       `json:"done_id" db:"done_id"`
	LastDoneAt *time.Time `json:"at" db:"at"`
}

// List :
func (repo *TaskRepository) List(option repository.ListOption) ([]model.Task, error) {
	sql := `
	SELECT
		t.*
		,done.id AS done_id
		,done.at
	FROM tasks t
	LEFT JOIN done_tasks done ON t.id = done.task_id
		AND NOT EXISTS (
			SELECT 1
			FROM done_tasks d
			WHERE t.id = d.task_id
			AND done.at < d.at
		)
	` + convertListOption(option)
	summaries := []TaskSummary{}
	if _, err := repo.Db.Select(&summaries, sql); err != nil {
		return nil, errors.WithStack(err)
	}

	tasks := make([]model.Task, len(summaries))
	taskIDs := make([]int, len(summaries))
	taskMap := make(map[int]*Task)
	for i, t := range summaries {
		task := &Task{
			TaskID:       t.TaskID,
			TaskName:     t.TaskName,
			TaskStartAt:  t.StartAt,
			TaskRuleType: t.TaskRuleType,
			TaskRule: &TaskRule{
				RuleType:      t.TaskRuleType,
				RulePeriods:   model.Periods{},
				RuleTimes:     model.Times{},
				RuleMonthDays: model.MonthDays{},
				RuleWeekdays:  model.Weekdays{},
				RuleDays:      model.Days{},
				RuleDates:     model.Dates{},
			},
		}
		if t.LastDoneID != nil {
			task.LastDoneTask = &DoneTask{
				DoneTaskID: *t.LastDoneID,
				TaskID:     task.TaskID,
				TaskName:   task.TaskName,
				DoneAt:     *t.LastDoneAt,
			}
		}
		tasks[i] = model.Task{TaskData: task}
		taskIDs[i] = task.TaskID
		taskMap[task.TaskID] = task
	}

	ruleLines, err := repo.ruleLineList(taskIDs...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, line := range ruleLines {
		task := taskMap[line.TaskID]
		task.TaskRule.addLine(line)
	}

	if option.Sort.By == repository.SortByTaskRemains {
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Deadline().Latest().Unix() < tasks[j].Deadline().Latest().Unix()
		})
	}

	return tasks, nil
}

// Create :
func (repo *TaskRepository) Create(transaction repository.Transaction, task *model.Task) error {
	trans := transaction.(*gorp.Transaction)

	rule := task.Rule()
	t := &Task{
		TaskName:     task.Name(),
		TaskStartAt:  task.StartAt(),
		TaskRuleType: rule.Type(),
	}
	if err := trans.Insert(t); err != nil {
		return errors.WithStack(err)
	}
	task.TaskData = t

	lines := repo.toRuleLines(t.TaskID, rule.TaskRuleData)
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

func (repo *TaskRepository) toRuleLines(taskID int, rule model.TaskRuleData) []TaskRuleLine {
	lines := []TaskRuleLine{}
	switch typ := rule.Type(); typ {
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

func (repo *TaskRepository) ruleLineList(taskIDs ...int) ([]TaskRuleLine, error) {
	lines := []TaskRuleLine{}
	if len(taskIDs) == 0 {
		return lines, nil
	}

	if _, err := repo.Db.Select(&lines, `
	SELECT *
	FROM task_rule_lines line
	WHERE line.task_id IN (:ids)
	`, map[string]interface{}{"ids": taskIDs}); err != nil {
		return nil, errors.WithStack(err)
	}
	return lines, nil
}

func (repo *TaskRepository) doneList(taskID int) ([]DoneTask, error) {
	dones := []DoneTask{}
	if _, err := repo.Db.Select(&dones, `
	SELECT *
	FROM done_tasks done
	WHERE done.task_id = ?
	`, taskID); err != nil {
		return nil, errors.WithStack(err)
	}
	return dones, nil
}

// Delete :
func (repo *TaskRepository) Delete(transaction repository.Transaction, task *model.Task) error {
	id := task.ID()
	dones, err := repo.doneList(id)
	if err != nil {
		return errors.WithStack(err)
	}
	ds := make([]interface{}, len(dones))
	for i, d := range dones {
		d := d
		ds[i] = &d
	}

	lines, err := repo.ruleLineList(id)
	if err != nil {
		return errors.WithStack(err)
	}
	ls := make([]interface{}, len(lines))
	for i, l := range lines {
		l := l
		ls[i] = &l
	}

	trans := transaction.(*gorp.Transaction)
	if _, err := trans.Delete(ds...); err != nil {
		return errors.WithStack(err)
	}

	if _, err := trans.Delete(ls...); err != nil {
		return errors.WithStack(err)
	}

	if _, err := trans.Delete(task.TaskData); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update :
func (repo *TaskRepository) Update(transaction repository.Transaction, task *model.Task) error {
	trans := transaction.(*gorp.Transaction)

	id := task.ID()
	rule := task.Rule()
	t := &Task{
		TaskID:       id,
		TaskName:     task.Name(),
		TaskStartAt:  task.StartAt(),
		TaskRuleType: rule.Type(),
	}

	if _, err := trans.Update(t); err != nil {
		return errors.WithStack(err)
	}
	task.TaskData = t

	{
		lines, err := repo.ruleLineList(id)
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
	}

	lines := repo.toRuleLines(id, rule.TaskRuleData)
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

// Done :
func (repo *TaskRepository) Done(task *model.Task, now time.Time) error {
	done := DoneTask{
		TaskID:   task.ID(),
		TaskName: task.Name(),
		DoneAt:   now,
	}
	if err := repo.Db.Insert(&done); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// One :
func (repo *TaskRepository) One(id int) (*model.Task, error) {
	var task Task
	err := repo.Db.SelectOne(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}

	ruleLines, err := repo.ruleLineList(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	task.TaskRule = &TaskRule{
		RuleType:      task.TaskRuleType,
		RulePeriods:   model.Periods{},
		RuleTimes:     model.Times{},
		RuleMonthDays: model.MonthDays{},
		RuleWeekdays:  model.Weekdays{},
		RuleDays:      model.Days{},
		RuleDates:     model.Dates{},
	}
	for _, line := range ruleLines {
		task.TaskRule.addLine(line)
	}

	return &model.Task{TaskData: &task}, nil
}

// Temporary :
func (repo *TaskRepository) Temporary(now time.Time) *model.Task {
	number := 1
	unit := model.PeriodUnitDay
	return &model.Task{TaskData: &Task{
		TaskName:     "name",
		TaskStartAt:  now,
		TaskRuleType: model.TaskRuleTypePeriodic,
		TaskRule: &TaskRule{
			RuleType:      model.TaskRuleTypePeriodic,
			RuleWeekdays:  model.Weekdays{},
			RuleDays:      model.Days{},
			RuleMonthDays: model.MonthDays{},
			RuleTimes:     model.Times{},
			RuleDates:     model.Dates{},
			RulePeriods: model.Periods{
				{
					PeriodData: TaskPeriod{
						PeriodNumber: &number,
						PeriodUnit:   &unit,
					},
				},
			},
		},
	}}
}

// From :
func (repo *TaskRepository) From(id int, reader io.Reader) (*model.Task, error) {
	task := &Task{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(task); err != nil {
		return nil, errors.WithStack(err)
	}
	task.TaskID = id
	return &model.Task{TaskData: task}, nil
}

// Task :
type Task struct {
	TaskID       int                `json:"id" db:"id, primarykey, autoincrement"`
	TaskName     string             `json:"name" db:"name, notnull" check:"notEmpty"`
	TaskStartAt  time.Time          `json:"startAt" db:"start_at, notnull"`
	TaskRuleType model.TaskRuleType `json:"ruleType" db:"rule_type, notnull" check:"taskRuleType"`

	LastDoneTask *DoneTask `json:"-" db:"-"`
	TaskRule     *TaskRule `json:"rule" db:"-"`
}

var _ model.TaskData = &Task{}

// ID :
func (task *Task) ID() int {
	return task.TaskID
}

// Name :
func (task *Task) Name() string {
	return task.TaskName
}

// Rule :
func (task *Task) Rule() *model.TaskRule {
	return &model.TaskRule{TaskRuleData: task.TaskRule}
}

// StartAt :
func (task *Task) StartAt() time.Time {
	return task.TaskStartAt
}

// LastDone :
func (task *Task) LastDone() *model.DoneTask {
	if task.LastDoneTask == nil {
		return nil
	}
	return &model.DoneTask{
		DoneTaskData: task.LastDoneTask,
	}
}

// DoneTask :
type DoneTask struct {
	DoneTaskID int       `json:"id" db:"id, primarykey, autoincrement"`
	TaskID     int       `json:"taskId" db:"task_id, notnull" foreign:"tasks(id)"`
	TaskName   string    `json:"name" db:"name, notnull" check:"notEmpty"`
	DoneAt     time.Time `json:"at" db:"at, notnull"`
}

// At :
func (done *DoneTask) At() time.Time {
	return done.DoneAt
}

var _ model.TaskRuleData = &TaskRule{}

// TaskRule :
type TaskRule struct {
	RuleType      model.TaskRuleType `json:"type"`
	RuleWeekdays  model.Weekdays     `json:"weekdays"`
	RuleDays      model.Days         `json:"days"`
	RuleMonthDays model.MonthDays    `json:"monthDays"`
	RuleTimes     model.Times        `json:"times"`
	RuleDates     model.Dates        `json:"dates"`
	RulePeriods   model.Periods      `json:"periods"`
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
	ID       int             `json:"id" db:"id, primarykey, autoincrement"`
	TaskID   int             `json:"taskId" db:"task_id, notnull" foreign:"tasks(id)"`
	Weekday  *time.Weekday   `json:"weekday" db:"weekday" check:"weekday"`
	Day      *model.Day      `json:"day" db:"day" check:"day"`
	MonthDay *model.MonthDay `json:"monthDay" db:"month_day"`
	Time     *time.Time      `json:"time" db:"time"`
	Date     *model.Date     `json:"date" db:"date"`
	TaskPeriod
}

var _ model.PeriodData = &TaskPeriod{}

// TaskPeriod :
type TaskPeriod struct {
	PeriodNumber *int              `json:"period_number" db:"period_number" check:"natural"`
	PeriodUnit   *model.PeriodUnit `json:"period_unit" db:"period_unit" check:"periodUnit"`
}

// Number :
func (period TaskPeriod) Number() int {
	return *period.PeriodNumber
}

// Unit :
func (period TaskPeriod) Unit() model.PeriodUnit {
	return *period.PeriodUnit
}
