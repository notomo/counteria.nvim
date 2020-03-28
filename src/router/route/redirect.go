package route

import (
	"fmt"
	"strconv"

	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"
)

// Redirector :
type Redirector struct {
	Vim *nvim.Nvim
}

// Do : redirect by route
func (re *Redirector) Do(r Route, params Params) error {
	path, err := r.BuildPath(params)
	if err != nil {
		return errors.WithStack(err)
	}

	// NOTE: avoid executing BufReadCmd

	var bufnr int
	if err := re.Vim.Call("bufnr", &bufnr, path); err != nil {
		return errors.WithStack(err)
	}

	var buf nvim.Buffer
	exists := bufnr != -1
	if !exists {
		b, err := re.Vim.CreateBuffer(false, true)
		if err != nil {
			return errors.WithStack(err)
		}
		buf = b
		if err := re.Vim.SetBufferName(buf, path); err != nil {
			return errors.WithStack(err)
		}
	}

	batch := re.Vim.NewBatch()
	batch.Command(fmt.Sprintf("edit %s", path))
	if err := batch.Execute(); err != nil {
		return errors.WithStack(err)
	}

	var unused interface{}
	if err := re.Vim.Call("counteria#read", unused, true, buf); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ToTasksOne :
func (re *Redirector) ToTasksOne(taskID int) error {
	return re.Do(TasksOne, Params{"taskId": strconv.Itoa(taskID)})
}
