package view

// Warn :
func (renderer *Renderer) Warn(msg string) error {
	var unused interface{}
	return renderer.Vim.Call("counteria#messenger#warn", unused, msg)
}
