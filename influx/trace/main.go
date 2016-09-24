// Package trace implements types to record influx traces. A trace encodes the
// state of an influx store as a state snapshot and a series of actions.  A
// different process may then rebuild a given influx state by loading the
// snapshot, then applying the actions in order.
package trace

import (
	"time"

	"github.com/nullstyle/go/influx"
)

// DefaultMaxSize is the maximum allowed size a snapshot can be unless otherwise
// specified.
const DefaultMaxSize = 1 * 1024 * 1024 // one megabyte

// ActionTrace represents the trace of a single action being dispatched against
// a store.
type ActionTrace struct {
	StartedAt  time.Time
	FinsihedAt time.Time
	Action     influx.Action
	Result     error
}

// Hook implement influx.AfterHook and influx.ErrorHook to record the trace
type Hook struct {
	Current    Snapshot
	MaxSize    uint64
	TargetSize uint64
	TargetAge  time.Duration

	next         ActionTrace
	snapshotting bool
}

var _ influx.BeforeHook = &Hook{}
var _ influx.AfterHook = &Hook{}
var _ influx.ErrorHook = &Hook{}
var _ influx.Named = &Hook{}

// Snapshot is a state snapshot.
type Snapshot struct {
	InitialState influx.Snapshot
	Dispatches   []ActionTrace
}
