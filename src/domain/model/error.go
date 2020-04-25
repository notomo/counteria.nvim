package model

import "fmt"

var (
	// ErrValidationRule :
	ErrValidationRule = fmt.Errorf("rule")
)

// ErrValidation :
type ErrValidation struct {
	Err     error
	Message string
}

func (e ErrValidation) Error() string {
	return fmt.Sprintf("%s: %s", e.Err.Error(), e.Message)
}

// NewErrValidation :
func NewErrValidation(err error, message string) error {
	return ErrValidation{
		Err:     err,
		Message: message,
	}
}
