package taskcmd

import (
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/notomo/counteria.nvim/src/lib"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Command :
type Command struct {
	Renderer   *view.BufferRenderer
	Buffer     *vimlib.BufferClient
	Redirector *route.Redirector
	Clock      lib.Clock

	TaskRepository     repository.TaskRepository
	TransactionFactory repository.TransactionFactory
}

// List :
func (cmd *Command) List() error {
	option := repository.ListOption{
		Sort: repository.Sort{
			By:    repository.SortByTaskRemains,
			Order: repository.SortOrderDesc,
		},
		Limit:  100,
		Offset: 0,
	}

	tasks, err := cmd.TaskRepository.List(option)
	if err != nil {
		return errors.WithStack(err)
	}

	now := cmd.Clock.Now()
	return cmd.Renderer.TaskList(tasks, now)
}

// Create :
func (cmd *Command) Create() error {
	reader, err := cmd.Buffer.Reader()
	if err != nil {
		return errors.WithStack(err)
	}

	var taskID int
	task, err := cmd.TaskRepository.From(taskID, reader)
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
	now := cmd.Clock.Now()
	task := cmd.TaskRepository.Temporary(now)
	return cmd.Renderer.OneNewTask(task)
}

// Delete :
func (cmd *Command) Delete(taskID int) error {
	task, err := cmd.TaskRepository.One(taskID)
	if err != nil {
		return errors.WithStack(err)
	}

	transaction, err := cmd.TransactionFactory.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := cmd.TaskRepository.Delete(transaction, task); err != nil {
		if err := transaction.Rollback(); err != nil {
			return errors.WithStack(err)
		}
		return errors.WithStack(err)
	}
	if err := transaction.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return cmd.Redirector.ToTasksList()
}

// Done :
func (cmd *Command) Done(taskID int) error {
	task, err := cmd.TaskRepository.One(taskID)
	if err != nil {
		return errors.WithStack(err)
	}

	now := cmd.Clock.Now()
	if err := cmd.TaskRepository.Done(task, now); err != nil {
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

	task, err := cmd.TaskRepository.From(taskID, reader)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.TaskRepository.Update(task); err != nil {
		return errors.WithStack(err)
	}

	return cmd.ShowOne(task.ID())
}
