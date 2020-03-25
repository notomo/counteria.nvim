package internal

import (
	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/router"
)

// Handler : rpc handler
type Handler struct {
	Router  *router.Router
	waiting chan struct{}
}

// NewHandler :
func NewHandler(router *router.Router) *Handler {
	return &Handler{
		Router:  router,
		waiting: make(chan struct{}),
	}
}

// Do : entry point for "do"
func (handler *Handler) Do(args []string) error {
	if err := handler.Router.Do(args); err != nil {
		return handler.Router.Error(err)
	}
	return nil
}

// Read : entry point for "read"
func (handler *Handler) Read(buf nvim.Buffer) error {
	if err := handler.Router.Read(buf); err != nil {
		return handler.Router.Error(err)
	}
	return nil
}

// Write : entry point for "write"
func (handler *Handler) Write(buf nvim.Buffer) error {
	if err := handler.Router.Write(buf); err != nil {
		return handler.Router.Error(err)
	}
	return nil
}

// NOTE: for testing

// StartWaiting : entry point for "startWaiting"
func (handler *Handler) StartWaiting() error {
	handler.waiting <- struct{}{}
	return nil
}

// Wait : entry point for "wait"
func (handler *Handler) Wait() error {
	<-handler.waiting
	return nil
}