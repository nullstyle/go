package influx

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Dispatch(t *testing.T) {
	state, store := baseTest(t)

	err := store.Dispatch(TestAction{Amount: 2})
	if assert.NoError(t, err) {
		assert.Equal(t, 2, state.Counter)
		assert.True(t, state.Child.Called, "child handler not called")
	}
}

func TestStore_Dispatch_AfterHook(t *testing.T) {
	_, store := baseTest(t)
	var called bool
	store.UseHooks(AfterDispatchFunc(func(ctx context.Context, action Action) error {
		called = true
		return nil
	}))

	// check after dispatch is called
	err := store.Dispatch(struct{}{})
	if assert.NoError(t, err) {
		assert.True(t, called, "after fn wasn't called")
	}
}

func TestStore_Dispatch_BeforeHook(t *testing.T) {
	state, store := baseTest(t)
	var called bool
	store.UseHooks(BeforeDispatchFunc(func(ctx context.Context, action Action) error {
		assert.Equal(t, 0, state.Counter, "before hook called _after_ state was manipulated")
		called = true
		return nil
	}))

	// check after dispatch is called
	err := store.Dispatch(struct{}{})
	if assert.NoError(t, err) {
		assert.True(t, called, "before fn wasn't called")
	}
}

func BenchmarkStore_Dispatch(b *testing.B) {
	cases := []struct {
		Name   string
		State  interface{}
		Action Action
	}{
		{"simple", &TestState{}, TestAction{Amount: 1}},
		{"bigger", &TestBiggerState{}, TestAction{Amount: 1}},
		// TODO: {"dynamic plan", &TestState{}, TestAction{Amount: 1}},
	}

	for _, kase := range cases {
		b.Run(kase.Name, func(b *testing.B) {

			store, err := New(kase.State)
			if err != nil {
				b.Errorf("error while creating store: %s", err)
				b.FailNow()
			}

			for i := 0; i < b.N; i++ {
				err := store.Dispatch(kase.Action)
				if err != nil {
					b.Errorf("error while dispatching: %s", err)
					b.Fail()
				}
			}
		})
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

func TestStore_UseHooks(t *testing.T) {
	_, store := baseTest(t)
	hook := &TestHook{}
	store.UseHooks(hook)
	assert.Len(t, store.hooks.after, 1)
	assert.Len(t, store.hooks.before, 1)
}
