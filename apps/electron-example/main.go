package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/electron"
	njs "github.com/nullstyle/go/gopherjs/js"
)

var win *js.Object

func main() {
	app := &electron.App{
		Name:         "go-electron-example",
		Version:      "0.1.0",
		WindowWidth:  640,
		WindowHeight: 480,
		OnReady:      ready,
	}

	electron.Start(app)
}

func ready(app *electron.App) {
	dirname := njs.Require("./node").Get("dirname").String()
	log.Println("in the ready callback:", dirname)
}
