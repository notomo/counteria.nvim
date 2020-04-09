package view

import (
	"bytes"
	"encoding/json"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// OneNewTask :
func (renderer *BufferRenderer) OneNewTask(task *model.Task) error {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(task.TaskData); err != nil {
		return errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))

	if err := renderer.Buffer.SetLines(
		lines[:len(lines)-1],
		renderer.Buffer.WithBufferType("acwrite"),
		renderer.Buffer.WithModifiable(true),
		renderer.Buffer.WithOpen(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// OneTask :
func (renderer *BufferRenderer) OneTask(task *model.Task) error {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(task.TaskData); err != nil {
		return errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))

	if err := renderer.Buffer.SetLines(
		lines[:len(lines)-1],
		renderer.Buffer.WithBufferType("acwrite"),
		renderer.Buffer.WithModifiable(true),
		renderer.Buffer.WithOpen(),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
