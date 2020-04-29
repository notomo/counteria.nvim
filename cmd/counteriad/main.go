package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"

	"github.com/notomo/counteria.nvim/cmd/counteriad/internal"
	"github.com/notomo/counteria.nvim/src/command"
	"github.com/notomo/counteria.nvim/src/datastore/sqliteimpl"
	"github.com/notomo/counteria.nvim/src/lib"
	"github.com/notomo/counteria.nvim/src/router"
	"github.com/notomo/counteria.nvim/src/router/route"
	"github.com/notomo/counteria.nvim/src/view"
	"github.com/notomo/counteria.nvim/src/vimlib"
)

var dataPath string

func init() {
	flag.StringVar(&dataPath, "data", "", "datastore file path")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(fmt.Sprintf("%+v", err))
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

	dep, err := sqliteimpl.Setup(
		sqliteimpl.WithDataPath(dataPath),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	bufClientFactory := &vimlib.BufferClientFactory{Vim: vim}
	handler := internal.NewHandler(
		router.New(
			vim,
			&command.RootCommand{
				Renderer:            &view.Renderer{Vim: vim},
				BufferClientFactory: bufClientFactory,
				Redirector:          &route.Redirector{Vim: vim, BufferClientFactory: bufClientFactory},
				Clock:               lib.NewClock(),
				Dep:                 dep,
			},
		),
	)

	handles := []struct {
		method string
		fn     interface{}
	}{
		{method: "do", fn: handler.Do},
		{method: "exec", fn: handler.Exec},
		// for buffer attach, detach
		{method: "nvim_buf_lines_event", fn: handler.BufferLinesEvent},
		{method: "nvim_buf_changedtick_event", fn: handler.BufferChangedtickEvent},
		{method: "nvim_buf_detach_event", fn: handler.BufferDetachEvent},
		// for testing
		{method: "start_waiting", fn: handler.StartWaiting},
		{method: "wait", fn: handler.Wait},
	}
	for _, h := range handles {
		if err := vim.RegisterHandler(h.method, h.fn); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := vim.Serve(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
