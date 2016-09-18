// Package electron implments functions for working with the electron main
// process.
package electron

import (
	"flag"
	"fmt"
	"log"
	"os"

	"runtime"

	"github.com/gopherjs/gopherjs/js"
	njs "github.com/nullstyle/go/gopherjs/js"
)

var (
	E = njs.Require("electron")
)

// App represents a go-electron app
type App struct {
	Name         string
	Version      string
	WindowWidth  int
	WindowHeight int
}

// On registers fn to be called whenever the global event provided is triggered.
// See available events at http://electron.atom.io/docs/api/app/#app.
func On(event string, fn func()) {
	E.Get("app").Call("on", event, fn)
}

// Start is the entrypoint for go-electron apps.
func Start(app *App) {
	flag.Parse()

	if *writePackageJSON {
		fmt.Printf(`{
      "name"    : "%s",
      "version" : "%s",
      "main"    : "main.js"
    }`, app.Name, app.Version)

		os.Exit(0)
	}

	// JS only section below
	if runtime.GOARCH != "js" {
		log.Fatal("abort: not running in gopherjs")
	}

	var win *js.Object

	On("activate", func() {
		if win == nil {
			win = createWindow(app.WindowWidth, app.WindowHeight)
		}
	})

	On("ready", func() {
		win = createWindow(app.WindowWidth, app.WindowHeight)
	})

	On("window-all-closed", func() {
		platform := js.Global.Get("process").Get("platform").String()
		if platform != "darwin" {
			E.Get("app").Call("quit")
		}
	})
}
