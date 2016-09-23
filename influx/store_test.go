package influx

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Dispatch(t *testing.T) {
	state, store := baseTest(t)

	err := store.Dispatch(context.Background(), TestAction{Amount: 2})
	if assert.NoError(t, err) {
		assert.Equal(t, 2, state.Counter)
		assert.True(t, state.Child.Called, "child handler not called")
	}
}

func TestStore_Dispatch_AfterHook(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		_, store := baseTest(t)
		var called bool
		store.UseHooks(AfterDispatchFunc(func(ctx context.Context, action Action) error {
			called = true
			return nil
		}))

		// check after dispatch is called
		err := store.Dispatch(context.Background(), struct{}{})
		if assert.NoError(t, err) {
			assert.True(t, called, "after fn wasn't called")
		}
	})

	t.Run("sad path: erroring hook", func(t *testing.T) {
		_, store := baseTest(t)
		store.UseHooks(AfterDispatchFunc(func(ctx context.Context, action Action) error {
			return errors.New("kaboom")
		}))

		// check after dispatch is called
		err := store.Dispatch(context.Background(), struct{}{})
		assert.Error(t, err)
	})
}

func TestStore_Dispatch_BeforeHook(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		state, store := baseTest(t)
		var called bool
		store.UseHooks(BeforeDispatchFunc(func(ctx context.Context, action Action) error {
			assert.Equal(t, 0, state.Counter, "before hook called _after_ state was manipulated")
			called = true
			return nil
		}))

		// check after dispatch is called
		err := store.Dispatch(context.Background(), struct{}{})
		if assert.NoError(t, err) {
			assert.True(t, called, "before fn wasn't called")
		}
	})

	t.Run("sad path: erroring hook", func(t *testing.T) {
		_, store := baseTest(t)
		store.UseHooks(BeforeDispatchFunc(func(ctx context.Context, action Action) error {
			return errors.New("kaboom")
		}))

		// check after dispatch is called
		err := store.Dispatch(context.Background(), struct{}{})
		if assert.Error(t, err) {

		}
	})
}

func TestStore_Dispatch_Context(t *testing.T) {
	_, store := baseTest(t)

	check := func(ctx context.Context, action Action) error {
		_, err := FromContext(ctx)
		assert.NoError(t, err)
		return nil
	}

	store.UseHooks(BeforeDispatchFunc(check))
	store.UseHooks(AfterDispatchFunc(check))

	// TODO: refactor test fixtures to allow us to report what phase a failure was
	// triggered in

	// ensure context has store available before, after dispatch
	err := store.Dispatch(context.Background(), "")
	assert.NoError(t, err)

	// ensure context is available during dispatch
	err = store.Dispatch(context.Background(), "check_store")
	assert.NoError(t, err)
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
				err := store.Dispatch(context.Background(), kase.Action)
				if err != nil {
					b.Errorf("error while dispatching: %s", err)
					b.Fail()
				}
			}
		})
	}
}

func TestStore_Dispatch_ErrorHook(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		_, store := baseTest(t)
		var called bool
		store.UseHooks(OnErrorFunc(func(ctx context.Context, action Action, e error) error {
			called = true
			return nil
		}))

		err := store.Dispatch(context.Background(), "")
		if assert.NoError(t, err) {
			assert.False(t, called, "error fn was called")
		}
	})

	t.Run("handler error", func(t *testing.T) {
		_, store := baseTest(t)
		var called bool
		store.UseHooks(OnErrorFunc(func(ctx context.Context, action Action, e error) error {
			called = true
			return nil
		}))

		err := store.Dispatch(context.Background(), "boom")
		assert.Error(t, err)
		assert.True(t, called, "error fn was not called")
	})

	t.Run("before hook error", func(t *testing.T) {
		_, store := baseTest(t)
		var called bool
		store.UseHooks(OnErrorFunc(func(ctx context.Context, action Action, e error) error {
			called = true
			return nil
		}))

		store.UseHooks(BeforeDispatchFunc(func(ctx context.Context, action Action) error {
			return errors.New("kaboom")
		}))

		err := store.Dispatch(context.Background(), "")
		assert.Error(t, err)
		// assert.True(t, called, "error fn was not called")
	})

	t.Run("error hook errors", func(t *testing.T) {
		_, store := baseTest(t)
		store.UseHooks(OnErrorFunc(func(ctx context.Context, action Action, e error) error {
			return errors.New("hook error")
		}))

		err := store.Dispatch(context.Background(), "boom")
		// TODO
		assert.Error(t, err)
	})
}

func TestStore_Dispatch_DoubleDispatch(t *testing.T) {
	// This is a known issue test that identifies and reproduces a bug in influx.
	// When using states that have embedded structs in them, the way we search for
	// handlers causes a second call to HandleAction to be made.

	var state struct {
		DispatchCount
		Correct DispatchCount
	}

	_, err := New(&state)
	require.NoError(t, err)
	// BROKEN: StateLoaded is being dispatched called twice
	assert.Equal(t, 2, state.DispatchCount.Value)

	// Working
	assert.Equal(t, 1, state.Correct.Value)

}

func TestStore_Go(t *testing.T) {
	_, store := baseTest(t)
	trigger := make(chan int)
	output := make(chan int, 1)

	// it runs in the background
	store.Go(func() {
		in := <-trigger
		output <- in
		close(output)
	})

	select {
	case <-output:
		assert.FailNow(t, "read from output channel before trigger")
	default:
		t.Log("good: still waiting on output")
	}

	trigger <- 4
	<-store.Done()
	select {
	case out := <-output:
		assert.Equal(t, 4, out)
	default:
		assert.FailNow(t, "output wasn't triggered")
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

func TestStore_LifecycleEvents(t *testing.T) {
	// wish we could use influxtext here.
	// just a simple smoke test for now.

	var state LifecycleTest

	store, err := New(&state)
	require.NoError(t, err)
	assert.True(t,
		state.LoadWasCalled,
		"StateLoaded wasn't called at store creation")
	assert.False(t,
		state.WillSaveWasCalled,
		"WillSaveWasCalled was called early")

	_, err = store.TakeSnapshot(context.Background())
	require.NoError(t, err)
	assert.True(t,
		state.WillSaveWasCalled,
		"StateWillSave wasn't called at save")

}

func TestStore_Save(t *testing.T) {
	state := TestState{
		Counter: 3,
	}
	store, err := New(&state)
	require.NoError(t, err)

	var out bytes.Buffer

	err = store.Save(context.Background(), &out)

	if assert.NoError(t, err) {
		var snap Snapshot
		var loaded TestState
		err = json.Unmarshal(out.Bytes(), &snap)
		require.NoError(t, err)
		err = json.Unmarshal(snap.State, &loaded)
		require.NoError(t, err)
		assert.Equal(t, 3, loaded.Counter)
	}

	// sad path: state fails to serialize as JSON
	chanState := make(chan int)
	store, err = New(&chanState)
	require.NoError(t, err)

	err = store.Save(context.Background(), ioutil.Discard)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "take snapshot failed")
	}

}

func TestStore_TakeSnapshot(t *testing.T) {
	state := TestState{
		Counter: 3,
	}
	store, err := New(&state)
	require.NoError(t, err)

	snap, err := store.TakeSnapshot(context.Background())
	if assert.NoError(t, err) {
		assert.True(t, snap.CreatedAt != time.Time{}, "CreatedAt isn't populated")

		var loaded TestState
		err := json.Unmarshal(snap.State, &loaded)
		assert.NoError(t, err)
		assert.Equal(t, 3, loaded.Counter)
	}

	// sad path: StateWillSave triggers an error
	var busted BreakAtWillSaveTest
	store, err = New(&busted)
	require.NoError(t, err)

	_, err = store.TakeSnapshot(context.Background())
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "StateWillSave dispatch failed")
	}

	// sad path: state fails to serialize as JSON
	chanState := make(chan int)
	store, err = New(&chanState)
	require.NoError(t, err)

	_, err = store.TakeSnapshot(context.Background())
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "encode state failed")
	}
}

func TestStore_Unwrap(t *testing.T) {
	state, store := baseTest(t)
	assert.Equal(t, state, store.Unwrap())
}

func TestStore_UseHooks(t *testing.T) {
	_, store := baseTest(t)
	hook := &TestHook{}
	store.UseHooks(hook)
	assert.Len(t, store.hooks.after, 1)
	assert.Len(t, store.hooks.before, 1)
	assert.Len(t, store.hooks.error, 1)
}
