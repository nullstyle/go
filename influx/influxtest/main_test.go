package influxtest

import (
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

	err := store.Dispatch(ActionInc)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, state.Counter)
	}

	err = store.Dispatch(ActionInc)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, state.Counter)
	}

	err = store.Dispatch(ActionDec)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, state.Counter)
	}

	err = store.Dispatch(ActionReset)
	if assert.NoError(t, err) {
		assert.Equal(t, 0, state.Counter)
	}

	err = store.Dispatch(ActionAfterHookError)
	assert.Error(t, err)

	err = store.Dispatch(ActionBeforeHookError)
	assert.Error(t, err)

	err = store.Dispatch(ActionDispatchError)
	assert.Error(t, err)

	err = store.Dispatch(struct{}{})
	assert.NoError(t, err, "test store is not extendable")
}
