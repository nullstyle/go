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
	req.lock.Lock()
	defer req.lock.Unlock()

	return (req.resp != nil || req.err != nil)
}

// Response returns the response state of the request, as of the latest
// dispatch.
func (req *Request) Response() (*http.Response, error) {
	req.lock.Lock()
	defer req.lock.Unlock()

	return req.resp, req.err
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

		state.lock.Lock()
		state.resp, state.err = action.Response()
		state.lock.Unlock()
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
	switch {
	case req.err != nil:
		state.State = req.err.Error()
	case req.resp != nil:
		state.State = req.resp.Status
	default:
		state.State = "pending"
	}

	enc, err := json.Marshal(state)
	if err != nil {
		return nil, errors.Wrap(err, "encode http state failed")
	}

	return enc, nil
}

func (req *Request) finish(resp *http.Response, err error) {
	req.lock.Lock()
	req.lock.Unlock()
	req.resp = resp
	req.err = err
}

var _ influx.Handler = &Request{}
var _ json.Marshaler = &Request{}
