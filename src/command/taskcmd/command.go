package taskcmd

import (
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/pkg/errors"
)

// Command :
type Command struct {
	Renderer       *view.BufferRenderer
	TaskRepository repository.TaskRepository
}

// List :
func (cmd *Command) List() error {
	tasks, err := cmd.TaskRepository.List()
	if err != nil {
		return errors.WithStack(err)
	}

	return cmd.Renderer.TaskList(tasks)
}

// Create :
func (cmd *Command) Create() error {
	task := model.Task{}
	return cmd.Renderer.OneNewTask(task)
}

// CreateForm :
func (cmd *Command) CreateForm() error {
	task := model.Task{}
	return cmd.Renderer.OneNewTask(task)
}
