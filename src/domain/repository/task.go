package repository

import (
	"time"

	"github.com/notomo/counteria.nvim/src/domain/model"
)

// TaskRepository :
type TaskRepository interface {
	List(ListOption) ([]model.Task, error)
	Create(Transaction, *model.Task) error
	Update(Transaction, *model.Task) error
	Delete(Transaction, *model.Task) error
	Done(Transaction, *model.Task, time.Time) error
	One(id int) (*model.Task, error)
	Temporary(now time.Time) *model.Task
}
