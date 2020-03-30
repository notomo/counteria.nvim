package router

import "fmt"

var (
	errInvalidRoute = fmt.Errorf("invalid route")
)

type routeErr struct {
	Err error
	Arg string
}

// Error :
func (rerr *routeErr) Error() string {
	return fmt.Sprintf("%s: %s", rerr.Err.Error(), rerr.Arg)
}

func newErr(typ error, arg string) error {
	return &routeErr{Err: typ, Arg: arg}
}
