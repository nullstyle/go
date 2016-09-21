package influx

import (
	"context"
	"reflect"
)

type contextKey int

// contextKeys holds the context keys for this package
var contextKeys struct {
	store contextKey
}

// handlert is the cached reflect.Type of the influx handler type. Used during
// dispatch to see if a given value or field implements Handler.
var handlert = reflect.TypeOf((*Handler)(nil)).Elem()

type afterFunc struct {
	Fn AfterFunc
}

func (hook *afterFunc) AfterDispatch(ctx context.Context, action Action) error {
	return hook.Fn(ctx, action)
}

type beforeFunc struct {
	Fn BeforeFunc
}

func (hook *beforeFunc) BeforeDispatch(
	ctx context.Context,
	action Action,
) error {
	return hook.Fn(ctx, action)
}

type errorFunc struct {
	Fn ErrorFunc
}

func (hook *errorFunc) DispatchError(
	ctx context.Context,
	action Action,
	e error,
) error {
	return hook.Fn(ctx, action, e)
}

var _ AfterHook = &afterFunc{}
var _ BeforeHook = &beforeFunc{}
var _ ErrorHook = &errorFunc{}
