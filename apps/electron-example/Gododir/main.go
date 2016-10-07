package main

import do "gopkg.in/godo.v2"

func tasks(p *do.Project) {
	p.Task("default", do.S{}, nil)

	// TODO
}

func main() {
	do.Godo(tasks)
}
