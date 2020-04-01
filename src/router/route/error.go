package route

import "fmt"

var (
	// ErrNotFound : 404
	ErrNotFound = fmt.Errorf("not found")
	// ErrInvalidAction : 400?
	ErrInvalidAction = fmt.Errorf("invalid action")
)

// Err :
type Err struct {
	Err error
	Arg string
}

// Error :
func (e *Err) Error() string {
	return fmt.Sprintf("%s: %s", e.Err.Error(), e.Arg)
}

// NewErrNotFound :
func NewErrNotFound(path string) error {
	return &Err{Err: ErrNotFound, Arg: path}
}

// NewErrInvalidAction :
func NewErrInvalidAction(action string) error {
	return &Err{Err: ErrInvalidAction, Arg: action}
}
