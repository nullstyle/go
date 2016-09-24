package main

import (
	"context"
	"log"

	"github.com/bep/gr"
	"github.com/nullstyle/go/influx"
)

// IncAction increments the counter
type IncAction struct {
	influx.Action
}

// DecAction decrements the counter
type DecAction struct {
	influx.Action
}

// ResetAction resets the counter to 0
type ResetAction struct {
	influx.Action
}

// State is the application's state
type State struct {
	influx.Handler
	Counter int
}

// UI represents the single ui component of this example.  In normal apps,
// you'd probably want to break this down into multiple components.
type UI struct {
	*gr.This

	Store *influx.Store
}

func main() {
	state := State{}
	store, err := influx.New(&state)
	if err != nil {
		log.Fatal(err)
	}

	ui := gr.New(&UI{
		Store: store,
	})

	renderer := influx.AfterDispatchFunc(func(
		ctx context.Context,
		action influx.Action,
	) error {
		log.Println("rendering")
		ui.Render("app", gr.Props{})
		return nil
	})

	store.UseHooks(renderer)
}
