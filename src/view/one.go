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

	buf := renderer.Buffer

	batch := renderer.Vim.NewBatch()
	batch.SetBufferLines(buf, 0, -1, true, lines[:len(lines)-1])
	batch.SetBufferOption(buf, "buftype", "acwrite")
	batch.SetBufferOption(buf, "modified", false)

	if err := batch.Execute(); err != nil {
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

	buf := renderer.Buffer

	batch := renderer.Vim.NewBatch()
	batch.SetBufferOption(buf, "modifiable", true)
	batch.SetBufferLines(buf, 0, -1, true, lines[:len(lines)-1])
	batch.SetBufferOption(buf, "buftype", "nofile")
	batch.SetBufferOption(buf, "modified", false)
	batch.SetBufferOption(buf, "modifiable", false)

	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
