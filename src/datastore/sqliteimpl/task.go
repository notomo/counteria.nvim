package sqliteimpl

import (
	"database/sql"
	"encoding/json"
	"io"
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
	TaskID     int        `json:"id" db:"id"`
	TaskName   string     `json:"name" db:"name"`
	LastDoneID *int       `json:"done_id" db:"done_id"`
	LastDoneAt *time.Time `json:"at" db:"at"`
	TaskPeriod
}

// List :
func (repo *TaskRepository) List() ([]model.Task, error) {
	summaries := []TaskSummary{}
	if _, err := repo.Db.Select(&summaries, `
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
	`); err != nil {
		return nil, errors.WithStack(err)
	}

	tasks := make([]model.Task, len(summaries))
	for i, t := range summaries {
		task := &Task{
			TaskID:     t.TaskID,
			TaskName:   t.TaskName,
			TaskPeriod: t.TaskPeriod,
		}
		if t.LastDoneID != nil {
			task.LastDone = &DoneTask{
				DoneTaskID: *t.LastDoneID,
				TaskID:     task.TaskID,
				TaskName:   task.TaskName,
				At:         *t.LastDoneAt,
			}
		}
		tasks[i] = task
	}
	return tasks, nil
}

// Create :
func (repo *TaskRepository) Create(task model.Task) error {
	if err := repo.Db.Insert(task); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete :
func (repo *TaskRepository) Delete(task model.Task) error {
	if _, err := repo.Db.Delete(task); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update :
func (repo *TaskRepository) Update(task model.Task) error {
	if _, err := repo.Db.Update(task); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Done :
func (repo *TaskRepository) Done(task model.Task, now time.Time) error {
	done := DoneTask{
		TaskID:   task.ID(),
		TaskName: task.Name(),
		At:       now,
	}
	if err := repo.Db.Insert(&done); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// One :
func (repo *TaskRepository) One(id int) (model.Task, error) {
	var task Task
	err := repo.Db.SelectOne(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, errors.WithStack(err)
	}
	return &task, nil
}

// Temporary :
func (repo *TaskRepository) Temporary() model.Task {
	return &Task{
		TaskName: "name",
		TaskPeriod: TaskPeriod{
			PeriodNumber: 1,
			PeriodUnit:   model.PeriodUnitDay,
		},
	}
}

// From :
func (repo *TaskRepository) From(id int, reader io.Reader) (model.Task, error) {
	task := &Task{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(task); err != nil {
		return nil, errors.WithStack(err)
	}
	task.TaskID = id
	return task, nil
}

// Task :
type Task struct {
	TaskID   int    `json:"id" db:"id"`
	TaskName string `json:"name" db:"name" check:"notEmpty"`
	TaskPeriod

	LastDone *DoneTask `json:"-" db:"-"`
}

var _ model.Task = &Task{}

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
	return task.TaskPeriod
}

// DoneAt :
func (task *Task) DoneAt() *time.Time {
	if task.LastDone == nil {
		return nil
	}
	return &task.LastDone.At
}

var _ model.Period = &TaskPeriod{}

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
	DoneTaskID int       `json:"id" db:"id"`
	TaskID     int       `json:"taskId" db:"task_id"`
	TaskName   string    `json:"name" db:"name" check:"notEmpty"`
	At         time.Time `json:"at" db:"at"`
}
