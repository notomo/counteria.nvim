package internal

import (
	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/router"
	"github.com/notomo/counteria.nvim/src/router/route"
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

// Exec : entry point for "exec"
func (handler *Handler) Exec(method route.Method, path string, bufnr int) error {
	buf := nvim.Buffer(bufnr)
	if err := handler.Router.Exec(method, path, buf); err != nil {
		return handler.Router.Error(err)
	}
	return nil
}

// BufferLinesEvent : entry point for "nvim_buf_lines_event"
func (handler *Handler) BufferLinesEvent(bufnr nvim.Buffer, changedtick, firstline, lastline int, linedata []string, more bool) error {
	return nil
}

// BufferChangedtickEvent : entry point for "nvim_buf_changedtick_event"
func (handler *Handler) BufferChangedtickEvent(bufnr nvim.Buffer, changedtick int) error {
	return nil
}

// BufferDetachEvent : entry point for "nvim_buf_detach_event"
func (handler *Handler) BufferDetachEvent(bufnr nvim.Buffer) error {
	return nil
}

// NOTE: for testing

// StartWaiting : entry point for "start_waiting"
func (handler *Handler) StartWaiting() error {
	handler.waiting <- struct{}{}
	return nil
}

// Wait : entry point for "wait"
func (handler *Handler) Wait() error {
	<-handler.waiting
	return nil
}
