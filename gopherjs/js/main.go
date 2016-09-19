// Package js implements gopherjs helper functions
package js

import (
	"runtime"

	gjs "github.com/gopherjs/gopherjs/js"
)

func Require(pkg string) *gjs.Object {
	if runtime.GOARCH != "js" {
		return nil
	}

	return gjs.Global.Call("require", pkg)
}
