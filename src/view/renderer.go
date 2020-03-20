package view

import "github.com/neovim/go-client/nvim"

// Renderer : render vim buffer
type Renderer struct {
	Vim *nvim.Nvim
}
