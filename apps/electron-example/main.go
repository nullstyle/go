package main

import (
	"log"

	"github.com/gopherjs/gopherjs/js"
	"github.com/nullstyle/go/electron"
	"github.com/nullstyle/go/gopherjs/module"
)

var win *js.Object

func main() {
	app := &electron.App{
		Name:         "go-electron-example",
		Version:      "0.1.0",
		WindowWidth:  1200,
		WindowHeight: 900,
		OnReady:      ready,
	}

	electron.Start(app)
}

func ready(app *electron.App) {
	dirname := module.Require("./node").Get("dirname").String()
	log.Println("in the ready callback:", dirname)
}
