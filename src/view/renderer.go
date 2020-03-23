package view

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/pkg/errors"
)

// Renderer :
type Renderer struct {
	Vim *nvim.Nvim
}

// BufferRenderer : for vim buffer
type BufferRenderer struct {
	Vim    *nvim.Nvim
	Buffer nvim.Buffer
}

// Buffer :
func (renderer *Renderer) Buffer(buf nvim.Buffer) *BufferRenderer {
	return &BufferRenderer{
		Vim:    renderer.Vim,
		Buffer: buf,
	}
}

// Decode :
func (renderer *BufferRenderer) Decode(result interface{}) error {
	b, err := renderer.Vim.BufferLines(renderer.Buffer, 0, -1, false)
	if err != nil {
		return errors.WithStack(err)
	}

	reader := bytes.NewReader(bytes.Join(b, nil))
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(result); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// SaveAndRedirect :
func (renderer *BufferRenderer) SaveAndRedirect(r route.Route, params route.Params) error {
	if err := renderer.Save(); err != nil {
		return errors.WithStack(err)
	}

	path, err := r.BuildPath(params)
	if err != nil {
		return errors.WithStack(err)
	}

	batch := renderer.Vim.NewBatch()
	batch.Command(fmt.Sprintf("edit %s", path))
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Save :
func (renderer *BufferRenderer) Save() error {
	batch := renderer.Vim.NewBatch()
	batch.SetBufferOption(renderer.Buffer, "modified", false)
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
