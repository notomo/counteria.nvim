package taskcmd

import "github.com/notomo/counteria.nvim/src/domain/model"

// Create :
func (cmd *Command) Create() error {
	task := model.Task{}
	return cmd.Renderer.OneTask(task)
}
