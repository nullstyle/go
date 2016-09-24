package influxtest

import (
	"context"
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/stretchr/testify/assert"
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

// AssertValue asserts the a response was seen and recorded, returning the found
// value if the assertion passes.
func (rt *ResponseTrap) AssertValue(t *testing.T) (interface{}, bool) {
	if !assert.NotNil(t, rt.SeenCtx, "request didn't complete") {
		return nil, false
	}
	val := rt.SeenCtx.Value(rt.Request)

	if !assert.NotNil(t, val, "no response found in context") {
		return nil, false
	}

	return val, true
}
