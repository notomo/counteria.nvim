package repository

import "github.com/notomo/counteria.nvim/src/domain/model"

// TaskRepository :
type TaskRepository interface {
	List() ([]model.Task, error)
	Create(*model.Task) error
	One(int) (*model.Task, error)
}
