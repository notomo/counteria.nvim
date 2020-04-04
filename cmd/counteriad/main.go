package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"

	"github.com/notomo/counteria.nvim/cmd/counteriad/internal"
	"github.com/notomo/counteria.nvim/src/command"
	"github.com/notomo/counteria.nvim/src/datastore/sqliteimpl"
	"github.com/notomo/counteria.nvim/src/router"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/notomo/counteria.nvim/src/vimlib"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	reader := os.Stdin
	writer := os.Stdout
	closer := os.Stdout

	f, err := os.OpenFile("/tmp/counteriad.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	log.SetOutput(f)
	logf := func(msg string, args ...interface{}) {
		log.Printf("%s: %s\n", msg, args)
	}

	vim, err := nvim.New(reader, writer, closer, logf)
	if err != nil {
		return errors.WithStack(err)
	}

	dep, err := sqliteimpl.Setup()
	if err != nil {
		return errors.WithStack(err)
	}

	handler := internal.NewHandler(
		router.New(
			vim,
			&command.RootCommand{
				Renderer:            &view.Renderer{Vim: vim},
				BufferClientFactory: &vimlib.BufferClientFactory{Vim: vim},
				Redirector:          &route.Redirector{Vim: vim},
				Dep:                 dep,
			},
		),
	)

	vim.RegisterHandler("do", handler.Do)
	vim.RegisterHandler("exec", handler.Exec)
	vim.RegisterHandler("execPath", handler.ExecPath)

	// for testing
	vim.RegisterHandler("startWaiting", handler.StartWaiting)
	vim.RegisterHandler("wait", handler.Wait)

	if err := vim.Serve(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
