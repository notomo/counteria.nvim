package vimlib

import "fmt"

var (
	// ErrNoState : it may not be plugin buffer
	ErrNoState = fmt.Errorf("no state")
)
