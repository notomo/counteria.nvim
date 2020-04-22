// +build tools

package main

import (
	_ "github.com/kyoh86/scopelint"
	_ "github.com/notomo/swityp/cmd/swityp"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
