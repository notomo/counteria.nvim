package router

import "fmt"

var (
	errInvalidRoute     = fmt.Errorf("invalid route")
	errInvalidReadPath  = fmt.Errorf("invalid read path")
	errInvalidWritePath = fmt.Errorf("invalid write path")
)

type routeErr struct {
	Err error
	Arg string
}

// Error :
func (rerr *routeErr) Error() string {
	return fmt.Sprintf("%s: %s", rerr.Err.Error(), rerr.Arg)
}
