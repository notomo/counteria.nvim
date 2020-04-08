package command

import (
	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/command/taskcmd"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/notomo/counteria.nvim/src/lib"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/notomo/counteria.nvim/src/vimlib"
)

// RootCommand :
type RootCommand struct {
	Renderer            *view.Renderer
	BufferClientFactory *vimlib.BufferClientFactory
	Redirector          *route.Redirector
	Clock               lib.Clock
	*domain.Dep
}

// TaskCmd :
func (root *RootCommand) TaskCmd(bufnr nvim.Buffer) *taskcmd.Command {
	client := root.BufferClientFactory.Get(bufnr)
	return &taskcmd.Command{
		Renderer:           root.Renderer.Buffer(client),
		Buffer:             client,
		Redirector:         root.Redirector,
		Clock:              root.Clock,
		TaskRepository:     root.TaskRepository,
		TransactionFactory: root.TransactionFactory,
	}
}
