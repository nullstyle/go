package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// IsDone returns true if the request has been completed
func (req *Request) IsDone() bool {
	if req.result == nil {
		return false
	}

	resp, err := req.result.get()

	return resp != nil || err != nil
}

// Result returns the response and error state of the request
func (req *Request) Result() (*http.Response, error) {
	if req.result == nil {
		return nil, nil
	}
	return req.result.get()
}

// HandleAction implements influx.Handler for the http system.
func (state *Request) HandleAction(
	ctx context.Context,
	action influx.Action,
) error {
	switch action := action.(type) {
	// a Response stored in the state tree copies the value of *Response
	// (disptached as an action against the store) when their IDs match.  While
	// this happens too late for parent components to respond to, this component
	// provides a cache of the response for later use.
	case *Request:
		if action.ID != state.ID {
			return nil
		}

		if state.IsDone() {
			return errors.New("duplicate http response seen")
		}

		// copy the result
		resp, err := action.Result()
		state.result = &result{}
		state.result.finish(resp, err)
	}

	return nil
}

// MarshalJSON implements json.Marshaler
func (req *Request) MarshalJSON() ([]byte, error) {
	// if protected, just return nil
	if !req.Request.Unprotected {
		return []byte("null"), nil
	}

	var state struct {
		ID    string
		State interface{}
	}
	resp, err := req.Result()
	switch {
	case err != nil:
		state.State = err.Error()
	case resp != nil:
		state.State = resp.Status
	default:
		state.State = "pending"
	}

	enc, err := json.Marshal(state)
	if err != nil {
		return nil, errors.Wrap(err, "encode http state failed")
	}

	return enc, nil
}

var _ influx.Handler = &Request{}
var _ json.Marshaler = &Request{}
