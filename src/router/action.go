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

	r, params, err := router.readRoutes.Match(state.Path)
	if err != nil {
		return newErr(errInvalidRoute, err.Error())
	}

	if err := router.Redirector.Do(r, params); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (router *Router) open(args []string) error {
	path := route.Schema + strings.Join(args, "")
	r, params, err := router.readRoutes.Match(path)
	if err != nil {
		return newErr(errInvalidRoute, err.Error())
	}

	if err := router.Redirector.Do(r, params); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
