package sqliteimpl

import (
	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// TaskRepository : impl
type TaskRepository struct {
	Db *gorp.DbMap
}

// List :
func (repo *TaskRepository) List() ([]model.Task, error) {
	return nil, nil
}

// Create :
func (repo *TaskRepository) Create(task *model.Task) error {
	if err := repo.Db.Insert(task); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// One :
func (repo *TaskRepository) One(id int) (*model.Task, error) {
	var task model.Task
	err := repo.Db.SelectOne(&task, "SELECT * FROM tasks WHERE ID=?", id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &task, nil
}
