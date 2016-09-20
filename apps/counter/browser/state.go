package main

import (
	"context"

	"github.com/nullstyle/go/influx"
)

// HandleAction implements influx.Handler
func (s *State) HandleAction(ctx context.Context, action influx.Action) error {
	switch action.(type) {
	case IncAction:
		s.Counter++
	case DecAction:
		s.Counter--
	case ResetAction:
		s.Counter = 0
	}

	return nil
}
