package influx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

// Dispatch applies the provided action to the store
func (store *Store) Dispatch(action Action) error {
	// TODO: think harder about concurrency protection
	store.lock.Lock()
	defer store.lock.Unlock()

	// create the dispatch context
	// TODO: probably want to add a configurable deadline or timeout here
	ctx := Context(context.Background(), store)

	// run the before dispatch hooks
	for i, hook := range store.hooks.before {
		err := hook.BeforeDispatch(ctx, action)
		if err != nil {
			store.errorHooks(ctx, store, err)
			return &HookError{Index: i, Hook: hook, Err: err}
		}
	}

	// TODO: run action middleware to mutate action
	// TODO: plan the dispatch

	// run the dispatch
	statev := reflect.ValueOf(store.state).Elem()
	err := store.dispatchValue(ctx, statev, action)
	if err != nil {
		store.errorHooks(ctx, store, err)
		return err
	}

	// run the after dispatch hooks
	for i, hook := range store.hooks.after {
		err := hook.AfterDispatch(ctx, action)
		if err != nil {
			store.errorHooks(ctx, store, err)
			return &HookError{Index: i, Hook: hook, Err: err}
		}
	}

	return nil
}

// Get loads dest with the current state
func (store *Store) Get(dest interface{}) error {
	destv := reflect.ValueOf(dest)
	statev := reflect.ValueOf(store.state)

	if statev.Type().AssignableTo(destv.Type()) {
		return errors.New("incorrect dest type")
	}

	destve := destv.Elem()
	if !destve.CanSet() {
		return errors.New("unsettable dest")
	}

	destve.Set(statev)
	return nil
}

// Save writes the state to w
func (store *Store) Save(w io.Writer) error {
	snapshot, err := store.TakeSnapshot()
	if err != nil {
		return errors.Wrap(err, "take snapshot failed")
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(&snapshot)
	if err != nil {
		return errors.Wrap(err, "encode snapshot failed")
	}

	return nil
}

// TakeSnapshot serializes the current state into a new snapshot
func (store *Store) TakeSnapshot() (Snapshot, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(store.state)
	if err != nil {
		return Snapshot{}, errors.Wrap(err, "encode state failed")
	}

	return Snapshot{
		CreatedAt: time.Now(),
		State:     json.RawMessage(buf.Bytes()),
	}, nil
}

// Unwrap returns the raw state value managed by this store.  Use with caution.
// TODO: add a vet check that ensures Unwrap is only used from test files.
func (store *Store) Unwrap() interface{} {
	store.lock.Lock()
	defer store.lock.Unlock()
	return store.state
}

// UseHooks adds a function that will be called after each successful
// dispatch of an action against the store
func (store *Store) UseHooks(hook Hook) {
	store.lock.Lock()
	defer store.lock.Unlock()

	before, ok := hook.(BeforeHook)
	if ok {
		store.hooks.before = append(store.hooks.before, before)
	}

	after, ok := hook.(AfterHook)
	if ok {
		store.hooks.after = append(store.hooks.after, after)
	}

	errhook, ok := hook.(ErrorHook)
	if ok {
		store.hooks.error = append(store.hooks.error, errhook)
	}
}

func (store *Store) dispatchValue(
	ctx context.Context,
	statev reflect.Value,
	action Action,
) error {
	// NOTE: should always be safe because the only possible statev values are
	// rooted at store.state, which should mean that any value is addressable.
	ptrv := statev.Addr()
	ptrt := ptrv.Type()

	if ptrt.Implements(handlert) {
		handler := ptrv.Interface().(Handler)

		err := handler.HandleAction(ctx, action)
		if err != nil {
			// TODO: wrap in a custom HandlerError type
			return &ActionError{
				Action: action,
				Store:  store,
				Err:    err,
			}
		}
	}

	// dispatch on children
	if statev.Kind() == reflect.Struct {
		for i := 0; i < statev.NumField(); i++ {
			child := statev.Field(i)
			err := store.dispatchValue(ctx, child, action)
			if err != nil {
				// NOTE: we wrap the error as usual since this is a recursive algorithm
				// TODO: figure out a better way to report where a dispatch failed
				return err
			}
		}
	}

	return nil
}

// errorHooks runs the error hooks registered with the store
func (store *Store) errorHooks(
	ctx context.Context,
	action Action,
	e error,
) {

	for i, hook := range store.hooks.error {
		err := hook.DispatchError(ctx, action, e)
		if err != nil {
			// NOTE(scott): I'm choosing to simply output an error triggered by an
			// error hook execution because IMO the original error is more important
			// to be bubbled up the stack and I don't want to introduce another error
			// type to the API.
			log.Print(&HookError{Index: i, Hook: hook, Err: err})
		}
	}

}