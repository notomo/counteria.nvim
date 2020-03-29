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
	Vim          *nvim.Nvim
	BufferClient *vimlib.BufferClient
}

// Buffer :
func (renderer *Renderer) Buffer(client *vimlib.BufferClient) *BufferRenderer {
	return &BufferRenderer{
		Vim:          renderer.Vim,
		BufferClient: client,
	}
}
