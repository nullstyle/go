// Package influxtest implements types that makes it easier to test your influx
// applications.
package influxtest

import (
	"testing"

	"reflect"

	"github.com/nullstyle/go/influx"
	"github.com/stretchr/testify/require"
)

var (
	ActionBeforeHookError = "error_in_before_hook"
	ActionAfterHookError  = "error_in_after_hook"
	ActionDispatchError   = "error_in_dispatch"
	ActionInc             = "inc"
	ActionDec             = "dec"
	ActionReset           = "reset"
)

// HookCase is a single test case for the behavior of a hook
type HookCase struct {
	Name    string
	Actions []influx.Action
	Hook    influx.Hook

	// Test is a func that takes a single parameter, which is called with the hook
	// for this test case after all provided actions are dispatched.
	Test interface{}

	BeforeInstall func(store *influx.Store)
}

// State is an influx state to aid with testing.
type State struct {
	Counter int
}

// Do runs the provided actions against the store, ensuring none of them cause
// dispatch to return an error.  Useful for setting up a test scenario.
func Do(t *testing.T, store *influx.Store, actions ...influx.Action) {
	for i, a := range actions {
		err := store.Dispatch(a)
		require.NoError(t, err, "action[%d] failed", i)
	}
}

// Hook runs the provided hook test cases.  This function runs through each
// provided case and: creates a new store, installs the hook, dispatches the
// actions, then calls the case's test function to assert on state.  Errors
// returned during action dispatch are ignored.
func Hook(t *testing.T, cases []HookCase) {
	for _, kase := range cases {
		store := New(t)

		if kase.BeforeInstall != nil {
			kase.BeforeInstall(store)
		}

		hook := kase.Hook
		store.UseHooks(hook)

		for _, action := range kase.Actions {
			store.Dispatch(action)
		}

		testv := reflect.ValueOf(kase.Test)
		if testv.Kind() != reflect.Func {
			require.Fail(t, "invalid Test field in hook case")
		}
		hookv := reflect.ValueOf(kase.Hook)
		testv.Call([]reflect.Value{hookv})
	}
}

// New returns a new instance of the influx test store
func New(t *testing.T, actions ...influx.Action) *influx.Store {
	var state State
	return new(t, &state, actions)
}

// NewFromState simplifies creating a store with a custom state in a testing
// context.
func NewFromState(
	t *testing.T,
	state interface{},
	actions ...influx.Action,
) *influx.Store {
	store, err := influx.New(state)
	require.NoError(t, err)

	Do(t, store, actions...)

	return store
}

func new(t *testing.T, state interface{}, actions []influx.Action) *influx.Store {
	store, err := influx.New(state)
	require.NoError(t, err)
	store.UseHooks(&afterHook{})
	store.UseHooks(&beforeHook{})

	Do(t, store, actions...)

	return store
}

var _ influx.Handler = &Handler{}
