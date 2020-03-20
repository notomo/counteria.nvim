package taskcmd

// List :
func (cmd *Command) List() error {
	tasks, err := cmd.TaskRepository.List()
	if err != nil {
		return err
	}

	return cmd.Renderer.TaskList(tasks)
}
