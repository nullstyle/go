package http

import (
	"net/http"
	"testing"

	"github.com/nullstyle/go/influx"
	"github.com/nullstyle/go/influx/influxtest"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {

	var state struct {
		Resp1 Response
		Resp2 Response
	}

	response := &Response{
		Request: influx.Request{
			ID: 1,
		},
		Resp: &http.Response{StatusCode: 200},
		Err:  nil,
	}

	state.Resp1.Request.ID = 1
	state.Resp2.Request.ID = 2

	influxtest.NewFromState(t, &state, response)

	if assert.True(t, state.Resp1.Done, "request 1 isn't done") {
		assert.Equal(t, response.Resp, state.Resp1.Resp)
		assert.Equal(t, response.Err, state.Resp1.Err)
	}

	if assert.False(t, state.Resp2.Done, "request 2 is mistakenly done") {
		assert.Nil(t, state.Resp2.Resp)
		assert.Nil(t, state.Resp2.Err)
	}
}
