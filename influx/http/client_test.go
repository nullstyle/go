package http

import (
	"context"
	"log"
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
	state.Result, err = client.Get(ctx, "https://google.com")

	if assert.NoError(t, err) {
		<-store.Done()

		if assert.True(t, state.Available(), "request isn't done yet") {
			resp, err := state.Get()
			log.Println("resp:", resp)
			log.Println("err:", err)
			if assert.NotNil(t, resp) {
				assert.Equal(t, 200, resp.StatusCode)
			}

			assert.NoError(t, err)
		}
	}
}

type clientTestState struct {
	Result
}
