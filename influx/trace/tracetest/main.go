// Package tracetest implements functions for testing an influx application
// using previously recorded execution traces.
package tracetest

import (
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/nullstyle/go/influx/trace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// JSONTrace represents a trace serialized as json
type JSONTrace struct {
	ID       string
	Snapshot []byte
}

type Trace interface {
	Name() string
	GetInitialState() (*influx.Store, error)
	GetDispatches() ([]trace.ActionTrace, error)
}

// MemoryTrace represents a trace in an in-memory trace.
type MemoryTrace struct {
	ID       string
	Snapshot trace.Snapshot
	Stores   StoreProvider
}

// ConfirmBroken ensures the trace runs exactly as it was recorded, with the
// same results.
func ConfirmBroken(t *testing.T, traces ...Trace) {
	for i, trace := range traces {
		store, err := trace.GetInitialState()
		require.NoError(t, err)

		dispatches, err := trace.GetDispatches()
		require.NoError(t, err)

		for _, dispatch := range dispatches {
			err := store.Dispatch(dispatch.Action)
			assert.Equal(t, dispatch.Result, err)
		}
	}
}

// ConfirmWorking ensures the trace runs without error.
func ConfirmWorking(t *testing.T, traces ...Trace) {
	for i, trace := range traces {
		store, err := trace.GetInitialState()
		require.NoError(t, err)

		dispatches, err := trace.GetDispatches()
		require.NoError(t, err)

		for _, dispatch := range dispatches {
			err := store.Dispatch(dispatch.Action)
			assert.NoError(t, err)
		}
	}
}

// TODO: add a run variant that actual waits wall clock time based upon the trace
