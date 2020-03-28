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
	Vim          *nvim.Nvim
	Root         *command.RootCommand
	BufferClient *vimlib.BufferClient
	Redirector   *route.Redirector
	readRoutes   route.Routes
	writeRoutes  route.Routes
}

// New :
func New(vim *nvim.Nvim, root *command.RootCommand) *Router {
	return &Router{
		Vim:          vim,
		Root:         root,
		BufferClient: root.Renderer.BufferClient,
		Redirector:   root.Renderer.Redirector,
		readRoutes: route.Routes{
			route.TasksNew,
			route.TasksOne,
			route.TasksList,
		},
		writeRoutes: route.Routes{
			route.TasksNew,
		},
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
		return &routeErr{errInvalidRoute, name}
	}

	if err := subRoute(args[1:]); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Read : from datastore to buffer
func (router *Router) Read(buf nvim.Buffer) error {
	path, err := router.Vim.BufferName(buf)
	if err != nil {
		return errors.WithStack(err)
	}

	r, params, err := router.readRoutes.Match(path)
	if err != nil {
		return &routeErr{errInvalidReadPath, err.Error()}
	}

	switch r {
	case route.TasksNew:
		return router.Root.TaskCmd(buf).CreateForm()
	case route.TasksOne:
		return router.Root.TaskCmd(buf).ShowOne(params["taskId"])
	case route.TasksList:
		return router.Root.TaskCmd(buf).List()
	}

	return &routeErr{errInvalidReadPath, path}
}

// Write : from buffer to datastore
func (router *Router) Write(buf nvim.Buffer) error {
	path, err := router.Vim.BufferName(buf)
	if err != nil {
		return errors.WithStack(err)
	}

	r, _, err := router.readRoutes.Match(path)
	if err != nil {
		return &routeErr{errInvalidWritePath, err.Error()}
	}

	switch r {
	case route.TasksNew:
		return router.Root.TaskCmd(buf).Create()
	}

	return &routeErr{errInvalidWritePath, path}
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
