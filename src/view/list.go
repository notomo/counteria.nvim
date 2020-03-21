package view

import (
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// TaskList :
func (renderer *BufferRenderer) TaskList(tasks []model.Task) error {
	lines := [][]byte{}
	for _, task := range tasks {
		lines = append(lines, []byte(task.Name))
	}

	buf := renderer.Buffer

	batch := renderer.Vim.NewBatch()
	batch.SetBufferLines(buf, 0, -1, false, lines)
	batch.SetBufferOption(buf, "modifiable", false)

	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
