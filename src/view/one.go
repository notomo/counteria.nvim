package view

import (
	"encoding/json"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/view/component"
	"github.com/pkg/errors"
)

// TaskFromForm :
func (renderer *BufferRenderer) TaskFromForm(taskID int) (*model.Task, error) {
	reader, err := renderer.Buffer.Reader()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var view component.TaskFormView
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&view); err != nil {
		return nil, errors.WithStack(err)
	}
	view.TaskID = taskID

	task := &model.Task{
		TaskData: &view,
	}
	if err := task.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	return task, nil
}

// OneTask : a task page
func (renderer *BufferRenderer) OneTask(task *model.Task) error {
	view := component.NewTaskForm(task)
	lines, err := view.Lines()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := renderer.Buffer.SetLines(
		lines,
		renderer.Buffer.WithBufferType("acwrite"),
		renderer.Buffer.WithModifiable(true),
		renderer.Buffer.WithOpen(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
