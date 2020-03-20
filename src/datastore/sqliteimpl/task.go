package sqliteimpl

import (
	"github.com/go-gorp/gorp"
	"github.com/notomo/counteria.nvim/src/domain/model"
)

// TaskRepository : impl
type TaskRepository struct {
	Db *gorp.DbMap
}

// List :
func (repo *TaskRepository) List() ([]model.Task, error) {
	return nil, nil
}
