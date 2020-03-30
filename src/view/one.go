package view

import (
	"bytes"
	"encoding/json"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// OneNewTask :
func (renderer *BufferRenderer) OneNewTask(task model.Task) error {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(task); err != nil {
		return errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))

	buffer := renderer.BufferClient
	if err := buffer.SetLines(
		lines[:len(lines)-1],
		buffer.WithBufferType("acwrite"),
		buffer.WithModifiable(true),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// OneTask :
func (renderer *BufferRenderer) OneTask(task model.Task) error {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(task); err != nil {
		return errors.WithStack(err)
	}
	lines := bytes.Split(b.Bytes(), []byte("\n"))

	buffer := renderer.BufferClient
	if err := buffer.SetLines(
		lines[:len(lines)-1],
		buffer.WithBufferType("acwrite"),
		buffer.WithModifiable(true),
	); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
