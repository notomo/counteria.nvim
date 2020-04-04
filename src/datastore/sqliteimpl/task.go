package sqliteimpl

import (
	"database/sql"
	"encoding/json"
	"io"

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

// List :
func (repo *TaskRepository) List() ([]model.Task, error) {
	tasks := []Task{}
	if _, err := repo.Db.Select(&tasks, "SELECT * FROM tasks"); err != nil {
		return nil, errors.WithStack(err)
	}
	ts := make([]model.Task, len(tasks))
	for i, t := range tasks {
		t := t
		ts[i] = &t
	}
	return ts, nil
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
func (repo *TaskRepository) From(reader io.Reader) (model.Task, error) {
	task := &Task{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(task); err != nil {
		return nil, errors.WithStack(err)
	}
	return task, nil
}

// Task :
type Task struct {
	TaskID   int    `json:"id" db:"id"`
	TaskName string `json:"name" db:"name" check:"notEmpty"`
	TaskPeriod
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
