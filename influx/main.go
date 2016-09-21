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

	"github.com/pkg/errors"
)

// Action represents a type that can affect an influx.State.
type Action interface {
}

type AfterFunc func(context.Context, Action) error

// AfterHook represents a hook function that handles the result of action an
// being dispatched against an influx store.
type AfterHook interface {
	AfterDispatch(context.Context, Action) error
}

// ActionError represents an error triggered during the application of an action
// to a store.
type ActionError struct {
	Action Action
	Store  *Store
	Cause  error
}

type BeforeFunc func(context.Context, Action) error

// BeforeHook represents a hook function that is triggered at the beginning of
// the dispatch process.
type BeforeHook interface {
	BeforeDispatch(context.Context, Action) error
}

type ErrorFunc func(context.Context, Action, error) error

// ErrorHook represents a hook function that is triggered whenever an error
// occurs durring dispatch.
type ErrorHook interface {
	DispatchError(ctx context.Context, action Action, e error) error
}

// Handler represents a portion of state that knows how to mutate itself in
// response to an action.
type Handler interface {
	HandleAction(ctx context.Context, action Action) error
}

// Hook represents a value that can plug in to the influx lifecycle.  A value
// provided as a Hook should implement one or more of the hook interfaces.  See
// the "Hooks" type for a list of available hooks.
type Hook interface{}

// HookError represents an error that occurred while running a hook function
type HookError struct {
	Index int
	Hook  Hook
	Err   error
}

// Hooks represents a type that implements all the possible influx hook
// interfaces.  It's never used directly, but defined here to document the
// possible implementations that can be provided to UseHooks()
type Hooks interface {
	AfterHook
	BeforeHook
	ErrorHook
}

// Named represent a value that know's its own name
type Named interface {
	Name() string
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
	lock  sync.Mutex
	state interface{}
	hooks struct {
		before []BeforeHook
		after  []AfterHook
		error  []ErrorHook
	}
}

// Context implements
func Context(parent context.Context, store *Store) context.Context {
	return context.WithValue(parent, &contextKeys.store, store)
}

// FromContext retrieves a *Store from the provided context
func FromContext(ctx context.Context) (*Store, error) {
	val := ctx.Value(&contextKeys.store)
	if val == nil {
		return nil, errors.New("no store in context")
	}

	ret, ok := val.(*Store)
	if !ok {
		return nil, errors.New("invalid store in context")
	}

	return ret, nil
}

// AfterDispatchFunc wraps the provided fn in a AfterHook implementation
func AfterDispatchFunc(fn AfterFunc) AfterHook {
	return &afterFunc{fn}
}

// BeforeDispatchFunc wraps the provided fn in a BeforeHook implementation
func BeforeDispatchFunc(fn BeforeFunc) BeforeHook {
	return &beforeFunc{fn}
}

// OnErrorFunc wraps the provided fn in a ErrorHook implementation
func OnErrorFunc(fn ErrorFunc) ErrorHook {
	return &errorFunc{fn}
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
