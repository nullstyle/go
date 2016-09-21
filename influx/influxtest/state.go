package influxtest

import (
	"context"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// HandleAction implements influx.Handler.
func (state *State) HandleAction(
	ctx context.Context,
	action influx.Action,
) error {

	switch action.(string) {
	case ActionInc:
		state.Counter++
	case ActionDec:
		state.Counter--
	case ActionReset:
		state.Counter = 0
	case ActionDispatchError:
		return errors.New("dispatch boom")
	}

	return nil
}

var _ influx.Handler = &State{}
