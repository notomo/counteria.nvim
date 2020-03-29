package router

import (
	"fmt"
	"log"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/command"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Router :
type Router struct {
	Vim                 *nvim.Nvim
	Root                *command.RootCommand
	BufferClientFactory *vimlib.BufferClientFactory
	Redirector          *route.Redirector
}

// New :
func New(vim *nvim.Nvim, root *command.RootCommand) *Router {
	return &Router{
		Vim:                 vim,
		Root:                root,
		BufferClientFactory: root.BufferClientFactory,
		Redirector:          root.Redirector,
	}
}

// Do : `:Counteria {args}`
func (router *Router) Do(args []string) error {
	if len(args) == 0 {
		return router.open(args)
	}

	var subRoute func(args []string) error
	switch name := args[0]; name {
	case "open":
		subRoute = router.open
	case "do":
		subRoute = router.do
	default:
		return newErr(errInvalidRoute, name)
	}

	if err := subRoute(args[1:]); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Request :
func (router *Router) Request(method route.Method, bufnr nvim.Buffer) error {
	path, err := router.Vim.BufferName(bufnr)
	if err != nil {
		return errors.WithStack(err)
	}

	r, params, err := route.All.Match(method, path)
	if err != nil {
		return newErr(errInvalidRoute, err.Error())
	}

	switch method {
	case route.MethodRead:
		switch r.Path {
		case route.TasksNew.Path:
			return router.Root.TaskCmd(bufnr).CreateForm()
		case route.TasksOne.Path:
			return router.Root.TaskCmd(bufnr).ShowOne(params.TaskID())
		case route.TasksList.Path:
			return router.Root.TaskCmd(bufnr).List()
		}
	case route.MethodWrite:
		switch r.Path {
		case route.TasksNew.Path:
			return router.Root.TaskCmd(bufnr).Create()
		}
	case route.MethodDelete:
		switch r.Path {
		case route.TasksOne.Path:
			return router.Root.TaskCmd(bufnr).Delete(params.TaskID())
		}
	}

	return newErr(errInvalidRoute, path)
}

// Error :
func (router *Router) Error(err error) error {
	if _, ok := err.(*routeErr); ok {
		msg := fmt.Sprintf("[counteria] %s\n", err)
		return router.Vim.WriteErr(msg)
	}

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
