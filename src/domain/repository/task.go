package repository

import (
	"io"

	"github.com/notomo/counteria.nvim/src/domain/model"
)

// TaskRepository :
type TaskRepository interface {
	List() ([]model.Task, error)
	Create(model.Task) error
	Update(model.Task) error
	Delete(model.Task) error
	One(id int) (model.Task, error)
	Temporary() model.Task
	From(id int, reader io.Reader) (model.Task, error)
}
