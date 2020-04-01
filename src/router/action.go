package router

import (
	"strings"

	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

func (router *Router) do(args []string) error {
	state, err := router.BufferClientFactory.Current().LoadLineState()
	if err != nil {
		if errors.Cause(err) == vimlib.ErrNoState {
			return route.NewErrInvalidAction(strings.Join(args, " "))
		}
		return errors.WithStack(err)
	}

	method := route.MethodRead
	if len(args) != 0 {
		switch args[0] {
		case "delete":
			method = route.MethodDelete
		default:
			return route.NewErrInvalidAction(strings.Join(args, " "))
		}
	}

	if err := router.Redirector.ToPath(method, state.Path); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (router *Router) open(args []string) error {
	path := route.Schema + strings.Join(args, "")
	if err := router.Redirector.ToPath(route.MethodRead, path); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
