// Package influx implements functions for building "unidirectional data flow"
// applications, specifically inspired by redux.
//
// An application built using this package centralizes its state into a single
// value (called a "store") and accepts modifications to the store by applying
// _actions_.
package influx

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"
)

// Action represents a type that can affect an influx.State.
type Action interface {
}

// AfterFunc represents a function that handles the result of action an
// being dispatched against an influx store.
type AfterFunc func(store *Store) error

// ActionError represents an error triggered during the application of an action
// to a store.
type ActionError struct {
	Action Action
	Store  *Store
	Cause  error
}

// Handler represents a portion of state that knows how to mutate itself in
// response to an action.
type Handler interface {
	HandleAction(ctx context.Context, action Action) error
}

// Snapshot is a snapshot of a store's state
type Snapshot struct {
	CreatedAt time.Time
	State     json.RawMessage
}

// Store contains a raw go value and manages the modification of that value.
// Actions are applied to the state one at a time. Avoid using circular
// references in your state object, as infinite recursion will occur.
type Store struct {
	lock     sync.Mutex
	state    interface{}
	afterFns []AfterFunc
	lastPlan []Handler
}

// New wraps the provided state in a new store
func New(state interface{}) (*Store, error) {
	statev := reflect.ValueOf(state)
	if statev.Kind() != reflect.Ptr {
		panic("state must be a pointer")
	}

	store := &Store{
		state: state,
	}

	// TODO: get initial action plan, cache it

	return store, nil
}
