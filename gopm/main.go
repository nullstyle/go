package gopm

import (
	"github.com/gopherjs/gopherjs/js"
)

// Require returns the exports object bound to the provided package name.  Will
// panic if not found.
func Require(pkg string) *js.Object {
	mods := js.Global.Get("gopm_modules")
	if mods == nil {
		panic("gopm_modules not found")
	}

	exports := mods.Get(pkg)
	if exports == nil {
		panic("package not found")
	}

	return exports
}
