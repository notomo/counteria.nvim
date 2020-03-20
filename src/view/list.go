package view

import (
	"github.com/notomo/counteria.nvim/src/domain/model"
)

// TaskList :
func (renderer *Renderer) TaskList(tasks []model.Task) error {
	buf, err := renderer.Vim.CreateBuffer(false, true)
	if err != nil {
		return err
	}

	lines := [][]byte{}
	for _, task := range tasks {
		lines = append(lines, []byte(task.Name))
	}

	batch := renderer.Vim.NewBatch()
	batch.SetBufferLines(buf, 0, -1, false, lines)
	batch.Command("tabedit")
	batch.Command("buffer " + string(buf))
	if err := batch.Execute(); err != nil {
		return err
	}

	return nil
}
