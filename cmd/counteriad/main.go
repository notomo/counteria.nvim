package main

import (
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neovim/go-client/nvim"

	"github.com/notomo/counteria.nvim/src/command"
	"github.com/notomo/counteria.nvim/src/datastore/sqliteimpl"
	"github.com/notomo/counteria.nvim/src/view"
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
		return err
	}
	defer f.Close()
	log.SetOutput(f)
	logf := func(msg string, args ...interface{}) {
		log.Printf("%s: %s\n", msg, args)
	}

	vim, err := nvim.New(reader, writer, closer, logf)
	if err != nil {
		return err
	}

	dep, err := sqliteimpl.Setup()
	if err != nil {
		return err
	}

	cmd := command.New(
		vim,
		&view.Renderer{Vim: vim},
		dep,
	)

	vim.RegisterHandler("do", cmd.Do)

	// for testing
	vim.RegisterHandler("wait", cmd.Wait)
	vim.RegisterHandler("startWaiting", cmd.StartWaiting)

	return vim.Serve()
}
