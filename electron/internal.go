package electron

import (
	"flag"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/gopherjs/module"
)

var (
	writePackageJSON = flag.Bool("writePackageJSON", false, "output package.json and exit")
)

func createWindow(width, height int) *js.Object {
	win := E.Get("BrowserWindow").New(map[string]int{
		"width":  width,
		"height": height,
	})

	dir := module.Require("./node").Get("dirname").String()
	path := fmt.Sprintf(`file://%s/index.html`, dir)
	// path := `file://index.html`

	win.Call("loadURL", path)

	// TODO: how should I best conditionalize the dev tools opening?
	// win.Get("webContents").Call("openDevTools")

	return win
}
