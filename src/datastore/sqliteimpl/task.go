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
	TaskID   int       `json:"id" db:"id"`
	TaskName string    `json:"name" db:"name"`
	StartAt  time.Time `json:"startAt" db:"start_at"`
	TaskPeriod

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
	for i, t := range summaries {
		task := &Task{
			TaskID:      t.TaskID,
			TaskName:    t.TaskName,
			TaskStartAt: t.StartAt,
			TaskPeriod:  t.TaskPeriod,
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
	}

	if option.Sort.By == repository.SortByTaskRemains {
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].LimitAt().Unix() < tasks[j].LimitAt().Unix()
		})
	}

	return tasks, nil
}

// Create :
func (repo *TaskRepository) Create(task *model.Task) error {
	if err := repo.Db.Insert(task.TaskData); err != nil {
		return errors.WithStack(err)
	}
	return nil
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
	dones, err := repo.doneList(task.ID())
	if err != nil {
		return errors.WithStack(err)
	}
	ds := make([]interface{}, len(dones))
	for i, d := range dones {
		d := d
		ds[i] = &d
	}

	trans := transaction.(*gorp.Transaction)
	if _, err := trans.Delete(ds...); err != nil {
		return errors.WithStack(err)
	}

	if _, err := trans.Delete(task.TaskData); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update :
func (repo *TaskRepository) Update(task *model.Task) error {
	if _, err := repo.Db.Update(task.TaskData); err != nil {
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
	return &model.Task{TaskData: &task}, nil
}

// Temporary :
func (repo *TaskRepository) Temporary(now time.Time) *model.Task {
	return &model.Task{TaskData: &Task{
		TaskName:    "name",
		TaskStartAt: now,
		TaskPeriod: TaskPeriod{
			PeriodNumber: 1,
			PeriodUnit:   model.PeriodUnitDay,
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
	TaskID      int       `json:"id" db:"id, primarykey, autoincrement"`
	TaskName    string    `json:"name" db:"name" check:"notEmpty"`
	TaskStartAt time.Time `json:"startAt" db:"start_at"`
	TaskPeriod

	LastDoneTask *DoneTask `json:"-" db:"-"`
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

// Period :
func (task *Task) Period() model.Period {
	return model.Period{PeriodData: task.TaskPeriod}
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

var _ model.PeriodData = &TaskPeriod{}

// TaskPeriod :
type TaskPeriod struct {
	PeriodNumber int              `json:"period_number" db:"period_number" check:"natural"`
	PeriodUnit   model.PeriodUnit `json:"period_unit" db:"period_unit" check:"periodUnit"`
}

// Number :
func (period TaskPeriod) Number() int {
	return period.PeriodNumber
}

// Unit :
func (period TaskPeriod) Unit() model.PeriodUnit {
	return period.PeriodUnit
}

// DoneTask :
type DoneTask struct {
	DoneTaskID int       `json:"id" db:"id, primarykey, autoincrement"`
	TaskID     int       `json:"taskId" db:"task_id" foreign:"tasks(id)"`
	TaskName   string    `json:"name" db:"name" check:"notEmpty"`
	DoneAt     time.Time `json:"at" db:"at"`
}

// At :
func (done *DoneTask) At() time.Time {
	return done.DoneAt
}
