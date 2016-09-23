package influxtest

import (
	"context"

	"github.com/nullstyle/go/influx"
)

// HandleAction implements influx.Handler for ResposneTrap. It causes the
// context of a dispatch to be recorded, provided the action is the request that
// the response trap has stored
func (rt *ResponseTrap) HandleAction(
	ctx context.Context,
	action influx.Action,
) error {
	switch action {
	case rt.Request:
		rt.SeenCtx = ctx
	}

	return nil
}
