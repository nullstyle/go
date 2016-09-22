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
		Req1 Request
		Req2 Request
	}

	req := &Request{
		Request: influx.Request{ID: 1},
		resp:    &http.Response{StatusCode: 200},
		err:     nil,
	}

	state.Req1.Request.ID = 1
	state.Req2.Request.ID = 2

	influxtest.NewFromState(t, &state, req)

	if assert.True(t, state.Req1.IsDone(), "request 1 isn't done") {
		resp, err := state.Req1.Response()
		assert.Equal(t, resp, state.Req1.resp)
		assert.Equal(t, err, state.Req1.err)
	}

	if assert.False(t, state.Req2.IsDone(), "request 2 is mistakenly done") {
		resp, err := state.Req2.Response()
		assert.Nil(t, resp)
		assert.Nil(t, err)
	}
}
