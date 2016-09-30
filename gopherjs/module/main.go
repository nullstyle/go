// Package module implements functions for import js modules
package module

import (
	"runtime"

	"github.com/gopherjs/gopherjs/js"
)

// Require imports module, returning its exports.  Defaults to using a global
// "require" method, but will fall back to a gopm packaged bundle if require
// isn't defined.
func Require(module string) *js.Object {
	if runtime.GOARCH != "js" {
		return nil
	}

	switch {
	case js.Global.Get("require") != nil:
		return js.Global.Call("require", module)
	case js.Global.Get("gopm_modules") != nil:
		return js.Global.Get("gopm_modules").Get(module)
	default:
		return nil
	}
}
