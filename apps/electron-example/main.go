package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/electron"
)

var win *js.Object

func main() {
	electron.On("activate", func() {
		if win == nil {
			createWindow()
		}
	})

	electron.On("ready", createWindow)
	electron.On("window-all-closed", func() {
		platform := js.Global.Get("process").Get("platform").String()
		if platform != "darwin" {
			electron.App.Call("quit")
		}
	})
}

func createWindow() {
	win = electron.BrowserWindow.New(map[string]int{
		"width":  400,
		"height": 400,
	})

	dir := js.Module.Get("__dirname").String()
	log.Println("hello", dir)
	// path := fmt.Sprintf(`file://%s/index.html`, dir)
	// path := `file://index.html`

	win.Call("loadURL", "https://github.com")
	// win.Get("webContents").Call("openDevTools")
}
