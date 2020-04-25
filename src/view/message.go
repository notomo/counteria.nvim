package view

// Warn :
func (renderer *Renderer) Warn(msg string) error {
	var unused interface{}
	return renderer.Vim.Call("counteria#messenger#warn", unused, msg)
}

// Err :
func (renderer *Renderer) Err(msg string) error {
	var unused interface{}
	return renderer.Vim.Call("counteria#messenger#error", unused, msg)
}
