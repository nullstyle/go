package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/electron"
)

var win *js.Object

func main() {
	app := &electron.App{
		Name:         "go-electron-example",
		Version:      "0.1.0",
		WindowWidth:  640,
		WindowHeight: 480,
	}

	electron.Start(app)
}
