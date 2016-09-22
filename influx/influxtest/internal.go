package influxtest

import (
	"context"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

type afterHook struct{}

var _ influx.AfterHook = &afterHook{}

func (h *afterHook) AfterDispatch(ctx context.Context, action influx.Action) error {
	switch action {
	case ActionAfterHookError:
		return errors.New("boom in after")
	}

	return nil
}

type beforeHook struct{}

var _ influx.BeforeHook = &beforeHook{}

func (h *beforeHook) BeforeDispatch(ctx context.Context, action influx.Action) error {
	switch action {
	case ActionBeforeHookError:
		return errors.New("boom in before")
	}

	return nil
}
