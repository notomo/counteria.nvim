package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/command/taskcmd"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/pkg/errors"
)

// Command : root command
type Command struct {
	Vim      *nvim.Nvim
	Renderer *view.Renderer
	*domain.Dep

	waiting chan struct{}
}

// New :
func New(vim *nvim.Nvim, renderer *view.Renderer, dep *domain.Dep) *Command {
	return &Command{
		Vim:      vim,
		Renderer: renderer,
		Dep:      dep,
		waiting:  make(chan struct{}),
	}
}

// Handle : `:Counteria {args...}`
func (cmd *Command) Handle(args []string) error {
	if len(args) == 0 {
		return errors.New("no args")
	}

	switch args[0] {
	case "task":
		cmd := taskcmd.Command{
			Renderer:       cmd.Renderer,
			TaskRepository: cmd.TaskRepository,
		}
		return cmd.Handle(args[1:])
	}

	return errors.Errorf("invalid args: %s", args)
}

// Do : entry point for rpc method "do"
func (cmd *Command) Do(args []string) error {
	if err := cmd.Handle(args); err != nil {
		return cmd.handleErr(err)
	}
	return nil
}

// StartWaiting : entry point for rpc method "startWaiting"
// NOTE: for testing
func (cmd *Command) StartWaiting() error {
	cmd.waiting <- struct{}{}
	return nil
}

// Wait : entry point for rpc method "wait"
// NOTE: for testing
func (cmd *Command) Wait() error {
	<-cmd.waiting
	return nil
}

func (cmd *Command) handleErr(err error) error {
	trace := fmt.Sprintf("%+v", err)
	lines := strings.Split(trace, "\n")
	msgs := []string{}
	for _, line := range lines {
		m := fmt.Sprintf("[countera] %s", strings.ReplaceAll(line, "\t", "    "))
		msgs = append(msgs, m)
	}
	msg := strings.Join(msgs, "\n")

	log.Println(msg)

	return cmd.Vim.WriteErr(msg)
}
