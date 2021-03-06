package router

import (
	"fmt"
	"log"
	"strings"

	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/command"
	"github.com/notomo/counteria.nvim/src/domain"
	"github.com/notomo/counteria.nvim/src/domain/model"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Router :
type Router struct {
	Vim                 *nvim.Nvim
	Root                *command.RootCommand
	BufferClientFactory *vimlib.BufferClientFactory
	Redirector          *route.Redirector
	Renderer            *view.Renderer
}

// New :
func New(vim *nvim.Nvim, root *command.RootCommand) *Router {
	return &Router{
		Vim:                 vim,
		Root:                root,
		BufferClientFactory: root.BufferClientFactory,
		Redirector:          root.Redirector,
		Renderer:            root.Renderer,
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
		return route.NewErrInvalidAction(name)
	}

	if err := subRoute(args[1:]); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Exec :
func (router *Router) Exec(method route.Method, path string, bufnr nvim.Buffer) error {
	var bufferPath string
	if path == "" {
		p, err := router.BufferClientFactory.Get(bufnr).Path()
		if err != nil {
			return errors.WithStack(err)
		}
		bufferPath = p
	} else {
		bufferPath = path
	}

	req, err := route.All.Match(method, bufferPath)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := router.exec(req, bufnr); err != nil {
		raw := errors.Cause(err)
		switch raw {
		case domain.ErrNotFound:
			return route.NewErrNotFound(bufferPath)
		}
		switch raw.(type) {
		case model.ErrValidation:
			return route.NewErrValidation(err)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (router *Router) exec(req route.Request, bufnr nvim.Buffer) error {
	path := req.Route.Path
	params := req.Params
	switch req.Method {
	case route.MethodRead:
		switch path {
		case route.TasksNew.Path:
			return router.Root.TaskCmd(bufnr).CreateForm()
		case route.TasksOne.Path:
			return router.Root.TaskCmd(bufnr).ShowOne(params.TaskID())
		case route.TasksList.Path:
			return router.Root.TaskCmd(bufnr).List()
		}
	case route.MethodWrite:
		switch path {
		case route.TasksNew.Path:
			return router.Root.TaskCmd(bufnr).Create()
		case route.TasksOne.Path:
			return router.Root.TaskCmd(bufnr).Update(params.TaskID())
		case route.TasksOneDone.Path:
			return router.Root.TaskCmd(bufnr).Done(params.TaskID())
		}
	case route.MethodDelete:
		switch path {
		case route.TasksOne.Path:
			return router.Root.TaskCmd(bufnr).Delete(params.TaskID())
		}
	}

	return route.NewErrNotFound(path)
}

// BufferLinesEvent :
func (router *Router) BufferLinesEvent(event vimlib.BufferLinesEvent) error {
	path, err := router.BufferClientFactory.Get(event.Bufnr).Path()
	if err != nil {
		return errors.WithStack(err)
	}

	req, err := route.Events.Match(route.MethodRead, path)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := router.event(req, event); err != nil {
		raw := errors.Cause(err)
		switch raw {
		case domain.ErrNotFound:
			return route.NewErrNotFound(path)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (router *Router) event(req route.Request, event vimlib.BufferLinesEvent) error {
	path := req.Route.Path
	switch path {
	// TODO
	}

	return route.NewErrNotFound(path)
}

// Error :
func (router *Router) Error(err error) error {
	if rerr, ok := errors.Cause(err).(*route.Err); ok {
		if rerr.IsWarn {
			return router.Renderer.Warn(err.Error())
		}
		return router.Renderer.Err(err.Error())
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

	return router.Vim.WritelnErr(msg)
}
