package route

import "fmt"

var (
	// ErrNotFound : 404
	ErrNotFound = fmt.Errorf("not found")
	// ErrInvalidAction : 400?
	ErrInvalidAction = fmt.Errorf("invalid action")
	// ErrValidation : create, update validation failed
	ErrValidation = fmt.Errorf("validation")
)

// Err :
type Err struct {
	Err    error
	Arg    string
	IsWarn bool
}

// Error :
func (e *Err) Error() string {
	return fmt.Sprintf("%s: %s", e.Err.Error(), e.Arg)
}

// NewErrNotFound :
func NewErrNotFound(path string) error {
	return &Err{Err: ErrNotFound, Arg: path, IsWarn: true}
}

// NewErrInvalidAction :
func NewErrInvalidAction(action string) error {
	return &Err{Err: ErrInvalidAction, Arg: action, IsWarn: true}
}

// NewErrValidation :
func NewErrValidation(err error) error {
	return &Err{Err: ErrValidation, Arg: err.Error()}
}
