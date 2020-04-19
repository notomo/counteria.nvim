package view

import (
	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/vimlib"
)

// Renderer :
type Renderer struct {
	Vim *nvim.Nvim
}

// BufferRenderer : for vim buffer
type BufferRenderer struct {
	*Renderer
	Vim    *nvim.Nvim
	Buffer *vimlib.BufferClient
}

// Buffer :
func (renderer *Renderer) Buffer(client *vimlib.BufferClient) *BufferRenderer {
	return &BufferRenderer{
		Renderer: renderer,
		Vim:      renderer.Vim,
		Buffer:   client,
	}
}
