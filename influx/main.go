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
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
)

//go:generate mockery -name Handler -output ./influxtest -case=underscore -outpkg=influxtest

// Action represents a action to be dispatched against a influx store. Handler
// attached to points of the state managed by an influx store interrogate the
// concrete type of an action to respond (or not respond) appropriately.
type Action interface {
}

type AfterFunc func(context.Context, Action) error

// AfterHook represents a hook function that handles the result of action an
// being dispatched against an influx store.
type AfterHook interface {
	AfterDispatch(context.Context, Action) error
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

// HandlerFunc represents a function that can respond to an action
type HandlerFunc func(ctx context.Context, action Action) error

// Hook represents a value that can plug in to the influx lifecycle.  A value
// provided as a Hook should implement one or more of the hook interfaces.  See
// the "Hooks" type for a list of available hooks.
type Hook interface{}

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

// TODO: I think this has been made obselete.  check into it after current commit
type Result struct {
	lock sync.Mutex
	req  Request
	ret  interface{}
	err  error
}

// Request is one of the fundamental methods of inter-component communication.
// It is used to build request/response style patterns of communication within
// influx.  For an example usage, see the influx/http package.
//
// NOTE: a request gets passed around between multiple go-routines.
type Request interface {

	// A request should be printable for debugging purposes
	fmt.Stringer
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
	tasks sync.WaitGroup

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

	err := store.init()
	if err != nil {
		return nil, errors.Wrap(err, "store initialization failed")
	}

	return store, nil
}

// NewRequest creates a new request
func NewRequest() Request {
	requestLock.Lock()
	id := nextRequest
	nextRequest++
	requestLock.Unlock()

	return requestID(id)
}
