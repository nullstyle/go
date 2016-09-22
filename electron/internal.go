package electron

import (
	"flag"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	njs "github.com/nullstyle/go/gopherjs/js"
)

var (
	writePackageJSON = flag.Bool("writePackageJSON", false, "output package.json and exit")
	writeAppDesc     = flag.Bool("writeAppDesc", false, "output the application json")
)

func createWindow(width, height int) *js.Object {
	win := E.Get("BrowserWindow").New(map[string]int{
		"width":  width,
		"height": height,
	})

	dir := njs.Require("./node").Get("dirname").String()
	path := fmt.Sprintf(`file://%s/index.html`, dir)

	win.Call("loadURL", path)

	// TODO: how should I best conditionalize the dev tools opening?
	// win.Get("webContents").Call("openDevTools")

	return win
}
