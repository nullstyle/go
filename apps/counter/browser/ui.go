package main

import (
	"fmt"
	"log"

	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
)

// Render renders the ui
func (ui UI) Render() gr.Component {
	var state State
	err := ui.Store.Get(&state)
	if err != nil {
		log.Println("error in render:", err)
	}
	message := fmt.Sprintf(" Click me! Number of clicks: %v", state.Counter)

	return el.Div(
		el.Button(
			gr.CSS("btn", "btn-lg", "btn-primary"),
			gr.Style("color", "orange"),
			gr.Text(message),
			evt.Click(ui.onClick)))
}

func (ui UI) onClick(event *gr.Event) {
	err := ui.Store.Dispatch(IncAction{})
	if err != nil {
		log.Println("error:", err)
	}
}
