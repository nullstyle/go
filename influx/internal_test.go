package influx

import (
	"context"
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
	switch action.(type) {
	case TestAction:
		state.Called = true
	}

	return nil
}

type TestState struct {
	Counter int
	Child   TestChild
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
