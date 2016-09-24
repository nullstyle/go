package influxtest

import (
	"context"
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	store := New(t, ActionInc, ActionInc, ActionInc)
	state := store.Unwrap().(*State)
	assert.Equal(t, 3, state.Counter)
}

func TestState(t *testing.T) {
	var state State
	store := new(t, &state, []influx.Action{})

	err := store.Dispatch(context.Background(), ActionInc)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, state.Counter)
	}

	err = store.Dispatch(context.Background(), ActionInc)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, state.Counter)
	}

	err = store.Dispatch(context.Background(), ActionDec)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, state.Counter)
	}

	err = store.Dispatch(context.Background(), ActionReset)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, state.Counter)
	}

	err = store.Dispatch(context.Background(), ActionAfterHookError)
	assert.Error(t, err)

	err = store.Dispatch(context.Background(), ActionBeforeHookError)
	assert.Error(t, err)

	err = store.Dispatch(context.Background(), ActionDispatchError)
	assert.Error(t, err)

	err = store.Dispatch(context.Background(), struct{}{})
	assert.NoError(t, err, "test store is not extendable")
}
