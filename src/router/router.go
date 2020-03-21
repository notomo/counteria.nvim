package router

import (
	"fmt"
	"log"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/command"
	"github.com/pkg/errors"
)

// Route :
type Route string

const schema = "counteria://"

var (
	// TaskNew :
	TaskNew = Route(schema + "task/new")
)

// Router :
type Router struct {
	Vim  *nvim.Nvim
	Root *command.RootCommand
}

// Do : `:Counteria {args}`
func (router *Router) Do(args []string) error {
	var route Route
	switch strings.Join(args, "/") {
	case "task/create":
		route = TaskNew
	default:
		return errors.Errorf("invalid args: %s", args)
	}

	batch := router.Vim.NewBatch()
	batch.Command(fmt.Sprintf("tabedit %s", route))
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Read : from datastore to buffer
func (router *Router) Read(buf nvim.Buffer) error {
	name, err := router.Vim.BufferName(buf)
	if err != nil {
		return errors.WithStack(err)
	}

	switch Route(name) {
	case TaskNew:
		return router.Root.TaskCmd(buf).CreateForm()
	}

	return errors.Errorf("invalid buffer name: %s", name)
}

// Write : from buffer to datastore
func (router *Router) Write(buf nvim.Buffer) error {
	name, err := router.Vim.BufferName(buf)
	if err != nil {
		return errors.WithStack(err)
	}

	switch Route(name) {
	case TaskNew:
		return router.Root.TaskCmd(buf).Create()
	}

	return errors.Errorf("invalid buffer name: %s", name)
}

// Error :
func (router *Router) Error(err error) error {
	trace := fmt.Sprintf("%+v", err)
	lines := strings.Split(trace, "\n")
	msgs := []string{}
	for _, line := range lines {
		m := fmt.Sprintf("[countera] %s", strings.ReplaceAll(line, "\t", "    "))
		msgs = append(msgs, m)
	}
	msg := strings.Join(msgs, "\n")

	log.Println(msg)

	return router.Vim.WriteErr(msg)
}
