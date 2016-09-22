package http

import (
	"context"
	"encoding/json"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// HandleAction implements influx.Handler for the http system.
func (resp *Response) HandleAction(
	ctx context.Context,
	action influx.Action,
) error {
	switch action := action.(type) {
	// a Response stored in the state tree copies the value of *Response
	// (disptached as an action against the store) when their IDs match.  While
	// this happens too late for parent components to respond to, this component
	// provides a cache of the response for later use.
	case *Response:
		if action.Request != resp.Request {
			return nil
		}

		if resp.Done {
			return errors.New("duplicate http response seen")
		}

		*resp = *action
		resp.Done = true
	}

	return nil
}

func (resp *Response) MarshalJSON() ([]byte, error) {
	// if protected, just return nil
	if !resp.Request.Unprotected {
		return []byte("null"), nil
	}

	var state struct {
		ID    string
		State interface{}
	}
	switch {
	case !resp.Done:
		state.State = "pending"
	case resp.Err != nil:
		state.State = resp.Err.Error()
	case resp.Resp != nil:
		state.State = resp.Resp.Status
	default:
		return nil, errors.New("invalid response state")
	}

	enc, err := json.Marshal(state)
	if err != nil {
		return nil, errors.Wrap(err, "encode http state failed")
	}

	return enc, nil
}
