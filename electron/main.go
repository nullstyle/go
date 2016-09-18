// Package electron implments functions for working with the electron main
// process.
package electron

import (
	"github.com/nullstyle/go/gopherjs/js"
)

var (
	E             = js.Require("electron")
	App           = E.Get("app")
	BrowserWindow = E.Get("BrowserWindow")
)

func On(event string, fn func()) {
	App.Call("on", event, fn)
}
