package main

import (
	"log"

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
		OnReady:      ready,
	}

	electron.Start(app)
}

func ready(app *electron.App) {
	log.Println("in the ready callback")
}
