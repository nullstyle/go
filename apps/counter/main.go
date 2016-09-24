package main

import (
	"github.com/nullstyle/go/electron"
)

func main() {
	app := &electron.App{
		Name:         "go-electron-counter",
		Version:      "0.1.0",
		WindowWidth:  640,
		WindowHeight: 480,
	}

	electron.Start(app)
}
