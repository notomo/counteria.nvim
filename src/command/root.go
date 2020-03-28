package command

import (
	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/command/taskcmd"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
)

// RootCommand :
type RootCommand struct {
	Renderer   *view.Renderer
	Redirector *route.Redirector
	*domain.Dep
}

// TaskCmd :
func (root *RootCommand) TaskCmd(buf nvim.Buffer) *taskcmd.Command {
	return &taskcmd.Command{
		Renderer:       root.Renderer.Buffer(buf),
		Redirector:     root.Redirector,
		TaskRepository: root.TaskRepository,
	}
}
