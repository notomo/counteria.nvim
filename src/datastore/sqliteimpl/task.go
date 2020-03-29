package sqliteimpl

import (
	"encoding/json"
	"io"

	"github.com/go-gorp/gorp"
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

// One :
func (repo *TaskRepository) One(id int) (model.Task, error) {
	var task Task
	err := repo.Db.SelectOne(&task, "SELECT * FROM tasks WHERE id=?", id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &task, nil
}

// Temporary :
func (repo *TaskRepository) Temporary() model.Task {
	return &Task{}
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
	TaskName string `json:"name" db:"name"`
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
