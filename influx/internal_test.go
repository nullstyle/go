package influx

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type BreakAtLoadTest struct {
}

func (state *BreakAtLoadTest) HandleAction(
	ctx context.Context,
	action Action,
) error {
	switch action {
	case StateLoaded:
		return errors.New("busted")
	}

	return errors.Errorf("unexpected message dispatched: %s", action)
}

type BreakAtWillSaveTest struct {
}

func (state *BreakAtWillSaveTest) HandleAction(
	ctx context.Context,
	action Action,
) error {
	switch action {
	case StateLoaded:
		return nil
	case StateWillSave:
		return errors.New("busted")
	}

	return errors.Errorf("unexpected message dispatched: %s", action)
}

// DispatchCount is a test component that records the count of dispatches that
// have called its HandleAction method.
type DispatchCount struct {
	Value int
}

func (state *DispatchCount) HandleAction(
	ctx context.Context,
	action Action,
) error {
	state.Value++
	return nil
}

type LifecycleTest struct {
	LoadWasCalled     bool
	WillSaveWasCalled bool
}

func (state *LifecycleTest) HandleAction(
	ctx context.Context,
	action Action,
) error {
	switch action {
	case StateLoaded:
		state.LoadWasCalled = true
		return nil
	case StateWillSave:
		state.WillSaveWasCalled = true
		return nil
	}

	return errors.Errorf("unexpected message dispatched: %s", action)
}

type TestAction struct {
	Amount int
}

type TestChild struct {
	Called bool
}

func (state *TestChild) HandleAction(ctx context.Context, action Action) error {
	switch action := action.(type) {
	case TestAction:
		state.Called = true
	case string:
		switch action {
		case "check_store":
			_, err := FromContext(ctx)
			return err
		case "boom":
			return errors.New("boom")
		}
	}

	return nil
}

type TestHook struct{}

var _ Hooks = &TestHook{}
var _ Named = &TestHook{}

func (hook *TestHook) Name() string {
	return "test-hook"
}

func (hook *TestHook) AfterDispatch(ctx context.Context, action Action) error {
	return nil
}
func (hook *TestHook) BeforeDispatch(ctx context.Context, action Action) error {
	return nil
}
func (hook *TestHook) DispatchError(ctx context.Context, action Action, err error) error {
	return nil
}

type TestState struct {
	Counter int
	Child   TestChild
}

type TestBiggerState struct {
	Child1 TestState
	Child2 TestState
	Child3 TestState
	Child4 TestState
	Child5 struct {
		Child6 TestState
	}
}

func (state *TestState) HandleAction(ctx context.Context, action Action) error {
	switch action := action.(type) {
	case TestAction:
		state.Counter += action.Amount
	}

	return nil
}

func baseTest(t *testing.T) (*TestState, *Store) {
	state := &TestState{}
	store, err := New(state)
	require.NoError(t, err)

	return state, store
}
