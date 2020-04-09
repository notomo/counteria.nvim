package route

import (
	"github.com/neovim/go-client/nvim"
	"github.com/notomo/counteria.nvim/src/vimlib"
	"github.com/pkg/errors"
)

// Redirector :
type Redirector struct {
	Vim                 *nvim.Nvim
	BufferClientFactory *vimlib.BufferClientFactory
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
		if err := re.BufferClientFactory.Current().SyncRequest(method.String(), path); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	buffer, err := re.BufferClientFactory.GetOrCreate(path)
	if err != nil {
		return errors.WithStack(err)
	}

	cursor, err := re.BufferClientFactory.Current().SaveCursor()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := buffer.SyncRequest(method.String(), path); err != nil {
		return errors.WithStack(err)
	}

	if err := cursor.Restore(); err != nil {
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
