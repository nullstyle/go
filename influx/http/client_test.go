package http

import (
	"context"
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/nullstyle/go/influx/influxtest"
	"github.com/stellar/go/support/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestClient_Get(t *testing.T) {

	http := httptest.NewClient()
	http.On("GET", "https://google.com").ReturnString(200, "hello")
	client := &Client{
		Raw: http,
	}

	// fails when not called from within a context that can read a store
	_, err := client.Get(context.Background(), "https://google.com")
	assert.EqualError(t, err, "get store failed: no store in context")

	var state clientTestState
	store := influxtest.NewFromState(t, &state)
	ctx := influx.Context(context.Background(), store)

	store.Do(func() {
		state.Resp.Request, err = client.Get(ctx, "https://google.com")
	})

	if assert.NoError(t, err) {
		// wait for all requests to complete (including any child dispatches they
		// trigger)
		<-store.Done()

		if val, ok := state.Resp.AssertValue(t); ok {
			if assert.IsType(t, val, result{}, val) {
				assert.Equal(t, 200, val.(result).Response.StatusCode)
			}
		}
	}
}

type clientTestState struct {
	Resp influxtest.ResponseTrap
}
