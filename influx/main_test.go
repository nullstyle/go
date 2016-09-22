package influx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	store := &Store{}
	ctx := Context(context.Background(), store)

	val := ctx.Value(&contextKeys.store)
	assert.Equal(t, store, val)
}

func TestFromContext(t *testing.T) {
	store := &Store{}
	ctx := context.WithValue(context.Background(), &contextKeys.store, store)

	// happy path
	found, err := FromContext(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, store, found)
	}

	// sad path
	found, err = FromContext(context.Background())
	assert.Error(t, err)

	// broken context: not a store at the store key
	ctx = context.WithValue(context.Background(), &contextKeys.store, 7)
	found, err = FromContext(ctx)
	assert.Error(t, err)

}

func TestNew(t *testing.T) {
	// works with pointers
	var (
		base TestState
		num  int
		pnum *int
		anon struct {
			Counter int
		}
	)

	New(&base)
	New(&num)
	New(&pnum)
	New(&anon)

	// panics when not a pointer
	assert.Panics(t, func() {
		New(TestState{})
	})

	// returns an error if state returns error on initial StateLoaded dispatch
	var busted BreakAtLoadTest
	_, err := New(&busted)
	assert.Error(t, err)
}
