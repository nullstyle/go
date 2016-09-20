package influx

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_AfterDispatch(t *testing.T) {
	_, store := baseTest(t)

	var called bool
	store.AfterDispatch(func(store *Store) error {
		called = true
		return nil
	})
	assert.Len(t, store.afterFns, 1)

	// check after dispatch is called
	err := store.Dispatch(struct{}{})
	if assert.NoError(t, err) {
		assert.True(t, called, "after fn wasn't called")
	}
}

func TestStore_Dispatch(t *testing.T) {
	state, store := baseTest(t)

	err := store.Dispatch(TestAction{Amount: 2})
	if assert.NoError(t, err) {
		assert.Equal(t, 2, state.Counter)
		assert.True(t, state.Child.Called, "child handler not called")
	}
}

func TestStore_Get(t *testing.T) {
	state, store := baseTest(t)
	state.Counter = 3

	var got *TestState
	err := store.Get(&got)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, got.Counter)
	}

	// unsettable
	err = store.Get(&TestState{})
	assert.Error(t, err)
}

func TestStore_Save(t *testing.T) {
	state := TestState{
		Counter: 3,
	}
	store, err := New(&state)
	require.NoError(t, err)

	var out bytes.Buffer
	err = store.Save(&out)
	if assert.NoError(t, err) {
		var snap Snapshot
		var loaded TestState
		err = json.Unmarshal(out.Bytes(), &snap)
		require.NoError(t, err)
		err = json.Unmarshal(snap.State, &loaded)
		require.NoError(t, err)

		assert.Equal(t, 3, loaded.Counter)
	}
}
