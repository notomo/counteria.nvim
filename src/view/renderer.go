package view

import (
	"github.com/neovim/go-client/nvim"
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
