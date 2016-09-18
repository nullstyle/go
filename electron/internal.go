package electron

import (
	"flag"
	"log"

	"github.com/gopherjs/gopherjs/js"
)

var (
	writePackageJSON = flag.Bool("writePackageJSON", false, "output package.json and exit")
)

func createWindow(width, height int) *js.Object {
	win := E.Get("BrowserWindow").New(map[string]int{
		"width":  width,
		"height": height,
	})

	dir := js.Module.Get("__dirname").String()
	log.Println("hello", dir)
	// path := fmt.Sprintf(`file://%s/index.html`, dir)
	// path := `file://index.html`

	win.Call("loadURL", "https://github.com")
	// win.Get("webContents").Call("openDevTools")

	return win
}
