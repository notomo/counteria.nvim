package view

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/pkg/errors"
)

// OneTask :
func (renderer *Renderer) OneTask(task model.Task) error {
	buf, err := renderer.Vim.CreateBuffer(false, true)
	if err != nil {
		return err
	}

	var bf bytes.Buffer
	encoder := json.NewEncoder(&bf)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(task); err != nil {
		return err
	}
	lines := bytes.Split(bf.Bytes(), []byte("\n"))

	batch := renderer.Vim.NewBatch()
	batch.SetBufferLines(buf, 0, -1, true, lines[:len(lines)-1])
	batch.SetBufferOption(buf, "buftype", "acwrite")
	batch.SetBufferOption(buf, "modified", false)
	batch.Command("tabedit")
	batch.Command(fmt.Sprintf("buffer %d", int(buf)))
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
