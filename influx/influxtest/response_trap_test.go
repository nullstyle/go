package influxtest

import (
	"context"
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseTrap_HandleAction(t *testing.T) {
	var state struct {
		R ResponseTrap
	}

	store := NewFromState(t, &state)
	req := influx.NewRequest()
	state.R.Request = req
	ctx := context.WithValue(context.Background(), req, 3)

	err := store.Dispatch(ctx, req)
	require.NoError(t, err)

	val := state.R.SeenCtx.Value(req)
	assert.Equal(t, 3, val)
}
