package trace

import (
	"context"
	"log"
	"time"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// AfterDispatch implements influx.AfterHook
func (hook *Hook) AfterDispatch(
	ctx context.Context,
	action influx.Action,
) error {
	hook.finish(nil)
	return nil
}

// BeforeDispatch implements influx.BeforeHook
func (hook *Hook) BeforeDispatch(
	ctx context.Context,
	action influx.Action,
) error {

	err := hook.manageSnapshot(ctx)
	if err != nil {
		return errors.Wrap(err, "manage snapshot failed")
	}

	hook.next = ActionTrace{
		StartedAt: time.Now(),
		Action:    action,
	}
	return nil
}

// DispatchError implements influx.ErrorHook
func (hook *Hook) DispatchError(
	ctx context.Context,
	action influx.Action,
	e error,
) error {
	hook.finish(e)
	return nil
}

// Name implements influx.Named
func (hook *Hook) Name() string {
	return "influx/trace"
}

func (hook *Hook) finish(err error) {
	hook.next.FinsihedAt = time.Now()
	hook.next.Result = err
	hook.Current.Dispatches = append(hook.Current.Dispatches, hook.next)
}

func (hook *Hook) manageSnapshot(ctx context.Context) error {
	if hook.snapshotting {
		return nil
	}

	// TODO: don't calculate size on every dispatch
	size, err := hook.Current.Size()
	if err != nil {
		return errors.Wrap(err, "failed to get trace size")
	}

	// the following block decides whether we should take a snapshot,
	// returning early if not, continuing through if a snapshot is needed
	switch {
	case hook.Current.Empty():
		// no-op
	case hook.Current.Age() > hook.TargetAge:
		// no-op
	case size > hook.MaxSize:
		return errors.New("snapshot max size exceeded")
	case size > hook.TargetSize:
		// no-op
	default:
		return nil
	}

	hook.snapshotting = true

	store, err := influx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "get store failed")
	}

	// schedule the snapshot for as soon as it is safe
	store.NextTick(ctx, func() {
		err = hook.Current.Checkpoint(store)
		if err != nil {
			// TODO: replace with store managed logging facility, when it exists
			log.Printf("WARN: checkpoint failed: %s", err)
		}
		hook.snapshotting = false
	})

	return nil
}
