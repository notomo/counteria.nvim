package sqliteimpl

import (
	"database/sql"
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

	RuleLines *TaskRuleLineRepository
	Dones     *DoneTaskRepository
}

var _ repository.TaskRepository = &TaskRepository{}

// TaskSummary :
type TaskSummary struct {
	TaskID       int                `db:"id"`
	TaskName     string             `db:"name"`
	StartAt      time.Time          `db:"start_at"`
	TaskRuleType model.TaskRuleType `db:"rule_type" check:"taskRuleType"`

	LastDoneID *int       `db:"done_id"`
	LastDoneAt *time.Time `db:"at"`
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

	ruleLines, err := repo.RuleLines.List(taskIDs...)
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

	if err := repo.RuleLines.Create(trans, t.TaskID, rule); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Delete :
func (repo *TaskRepository) Delete(transaction repository.Transaction, task *model.Task) error {
	trans := transaction.(*gorp.Transaction)

	if err := repo.Dones.Delete(trans, task); err != nil {
		return errors.WithStack(err)
	}

	if err := repo.RuleLines.Delete(trans, task.ID()); err != nil {
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

	rule := task.Rule()
	t := &Task{
		TaskID:       task.ID(),
		TaskName:     task.Name(),
		TaskStartAt:  task.StartAt(),
		TaskRuleType: rule.Type(),
	}

	if _, err := trans.Update(t); err != nil {
		return errors.WithStack(err)
	}
	task.TaskData = t

	if err := repo.RuleLines.Update(trans, task.ID(), rule); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Done :
func (repo *TaskRepository) Done(transaction repository.Transaction, task *model.Task, now time.Time) error {
	if err := repo.Dones.Create(transaction, task, now); err != nil {
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

	ruleLines, err := repo.RuleLines.List(id)
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

// Task :
type Task struct {
	TaskID       int                `db:"id, primarykey, autoincrement"`
	TaskName     string             `db:"name, notnull" check:"notEmpty"`
	TaskStartAt  time.Time          `db:"start_at, notnull"`
	TaskRuleType model.TaskRuleType `db:"rule_type, notnull" check:"taskRuleType"`

	LastDoneTask *DoneTask `db:"-"`
	TaskRule     *TaskRule `db:"-"`
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
