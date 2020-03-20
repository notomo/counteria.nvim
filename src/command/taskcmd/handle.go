package taskcmd

import (
	"github.com/notomo/counteria.nvim/src/domain/repository"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/pkg/errors"
)

// Command :
type Command struct {
	Renderer       *view.Renderer
	TaskRepository repository.TaskRepository
}

// Handle :
func (cmd *Command) Handle(args []string) error {
	if len(args) == 0 {
		return cmd.List()
	}

	switch args[0] {
	case "list":
		return cmd.List()
	case "create":
		return cmd.Create()
	}

	return errors.Errorf("invalid args: %s", args)
}
