package taskcmd

import (
	"strconv"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/pkg/errors"
)

// Command :
type Command struct {
	Renderer       *view.BufferRenderer
	Redirector     *route.Redirector
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
	var task model.Task
	if err := cmd.Renderer.Decode(&task); err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.TaskRepository.Create(&task); err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.Renderer.Save(); err != nil {
		return errors.WithStack(err)
	}

	return cmd.Redirector.ToTasksOne(task.ID)
}

// ShowOne :
func (cmd *Command) ShowOne(taskID string) error {
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return errors.WithStack(err)
	}

	task, err := cmd.TaskRepository.One(id)
	if err != nil {
		return errors.WithStack(err)
	}
	return cmd.Renderer.OneTask(*task)
}

// CreateForm :
func (cmd *Command) CreateForm() error {
	task := model.Task{}
	return cmd.Renderer.OneNewTask(task)
}
