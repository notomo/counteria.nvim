package route

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
	"github.com/pkg/errors"
)

// Redirector :
type Redirector struct {
	Vim *nvim.Nvim
}

// To : redirect by route
func (re *Redirector) To(method Method, r Route, params Params) error {
	path, err := r.BuildPath(params)
	if err != nil {
		return errors.WithStack(err)
	}
	return re.ToPath(method, path)
}

// ToPath : redirect by path
func (re *Redirector) ToPath(method Method, path string) error {
	if !method.Renderable() {
		var unused interface{}
		var buf nvim.Buffer
		if err := re.Vim.Call("counteria#request_path", unused, method, true, path, buf); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	// NOTE: avoid executing BufReadCmd

	var bufnr int
	if err := re.Vim.Call("bufnr", &bufnr, path); err != nil {
		return errors.WithStack(err)
	}

	buf := nvim.Buffer(bufnr)
	exists := bufnr != -1
	if !exists {
		b, err := re.Vim.CreateBuffer(false, true)
		if err != nil {
			return errors.WithStack(err)
		}
		buf = b
		bufnr = int(buf)
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
	if err := re.Vim.Call("counteria#request", unused, method, true, buf); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// ToTasksOne :
func (re *Redirector) ToTasksOne(taskID int) error {
	path, err := TasksOnePath(taskID)
	if err != nil {
		return errors.WithStack(err)
	}
	return re.ToPath(MethodRead, path)
}

// ToTasksList :
func (re *Redirector) ToTasksList() error {
	return re.To(MethodRead, TasksList, Params{})
}
