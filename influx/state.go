package influx

import (
	"context"
	"encoding/json"
	"io"
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
	ctx := context.TODO()

	// run the before dispatch hooks
	for _, hook := range store.hooks.before {
		err := hook.BeforeDispatch(ctx, action)
		if err != nil {
			return errors.Wrap(err, "after-dispatch failed")
		}
	}

	// TODO: run action middleware to mutate action
	// TODO: plan the dispatch

	// run the dispatch
	statev := reflect.ValueOf(store.state).Elem()
	err := store.dispatchValue(ctx, statev, action)
	if err != nil {
		return errors.Wrap(err, "dispatch-value failed")
	}

	// run the after dispatch hooks
	for _, hook := range store.hooks.after {
		err := hook.AfterDispatch(ctx, action)
		if err != nil {
			return errors.Wrap(err, "after-dispatch failed")
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
	snapshot := struct {
		CreatedAt time.Time
		State     interface{}
	}{
		CreatedAt: time.Now(),
		State:     store.state,
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(snapshot)
	if err != nil {
		return errors.Wrap(err, "encode snapshot failed")
	}

	return nil
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
				Cause:  err,
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