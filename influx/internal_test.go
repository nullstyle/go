package influx

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		case "boom":
			log.Println("in boom")
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
