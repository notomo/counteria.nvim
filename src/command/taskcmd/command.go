package taskcmd

import (
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Command :
type Command struct {
	Renderer       *view.BufferRenderer
	Buffer         *vimlib.BufferClient
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
	reader, err := cmd.Buffer.Reader()
	if err != nil {
		return errors.WithStack(err)
	}

	task, err := cmd.TaskRepository.From(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.TaskRepository.Create(task); err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.Buffer.Save(); err != nil {
		return errors.WithStack(err)
	}

	return cmd.Redirector.ToTasksOne(task.ID())
}

// ShowOne :
func (cmd *Command) ShowOne(taskID int) error {
	task, err := cmd.TaskRepository.One(taskID)
	if err != nil {
		return errors.WithStack(err)
	}
	return cmd.Renderer.OneTask(task)
}

// CreateForm :
func (cmd *Command) CreateForm() error {
	task := cmd.TaskRepository.Temporary()
	return cmd.Renderer.OneNewTask(task)
}

// Delete :
func (cmd *Command) Delete(taskID int) error {
	task, err := cmd.TaskRepository.One(taskID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.TaskRepository.Delete(task); err != nil {
		return errors.WithStack(err)
	}

	return cmd.Redirector.ToTasksList()
}

// Update :
func (cmd *Command) Update(taskID int) error {
	reader, err := cmd.Buffer.Reader()
	if err != nil {
		return errors.WithStack(err)
	}

	task, err := cmd.TaskRepository.From(reader)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.TaskRepository.Update(task); err != nil {
		return errors.WithStack(err)
	}

	return cmd.ShowOne(task.ID())
}
