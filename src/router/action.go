package router

import (
	"strings"

	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/pkg/errors"
)

func (router *Router) do(args []string) error {
	state, err := router.BufferClient.LoadLineState()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := router.Redirector.ToPath(state.Path); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (router *Router) open(args []string) error {
	path := route.Schema + strings.Join(args, "")
	if err := router.Redirector.ToPath(path); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
