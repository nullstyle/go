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
		IndexStylesheets: []string{
			"https://maxcdn.bootstrapcdn.com/bootswatch/3.3.6/darkly/bootstrap.min.css",
		},
		IndexScripts: []string{
			"https://fb.me/react-15.0.1.min.js",
			"https://fb.me/react-dom-15.0.1.min.js",
		},
	}

	electron.Start(app)
}
