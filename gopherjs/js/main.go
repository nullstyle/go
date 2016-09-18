// Package js implements gopherjs helper functions
package js

import (
	gjs "github.com/gopherjs/gopherjs/js"
)

func Require(pkg string) *gjs.Object {
	return gjs.Global.Call("require", pkg)
}
