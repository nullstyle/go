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
		result:  &result{},
	}
	req.result.finish(&http.Response{StatusCode: 200}, nil)

	state.Req1.Request.ID = 1
	state.Req2.Request.ID = 2

	influxtest.NewFromState(t, &state, req)

	if assert.True(t, state.Req1.IsDone(), "request 1 isn't done") {
		actual, aerr := state.Req1.Result()
		expected, eerr := state.Req1.Result()
		assert.Equal(t, expected, actual)
		assert.Equal(t, eerr, aerr)
	}

	if assert.False(t, state.Req2.IsDone(), "request 2 is mistakenly done") {
		resp, err := state.Req2.Result()
		assert.Nil(t, resp)
		assert.Nil(t, err)
	}
}
