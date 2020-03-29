package view

import (
	"sync"

	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Renderer :
type Renderer struct {
	Vim          *nvim.Nvim
	BufferClient *vimlib.BufferClient
	Redirector   *route.Redirector

	NsID         int
	getNamespace sync.Once
}

// BufferRenderer : for vim buffer
type BufferRenderer struct {
	Vim          *nvim.Nvim
	Buffer       nvim.Buffer
	BufferClient *vimlib.BufferClient
	NsID         int
	Redirector   *route.Redirector
}

// Buffer :
func (renderer *Renderer) Buffer(buf nvim.Buffer) *BufferRenderer {
	renderer.getNamespace.Do(func() {
		ns, err := renderer.Vim.CreateNamespace("counteria")
		if err != nil {
			panic(err)
		}
		renderer.NsID = ns
	})

	return &BufferRenderer{
		Vim:          renderer.Vim,
		Buffer:       buf,
		BufferClient: renderer.BufferClient,
		NsID:         renderer.NsID,
		Redirector:   renderer.Redirector,
	}
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
