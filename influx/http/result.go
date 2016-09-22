package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/nullstyle/go/influx"
)

func (r *Result) Get() (*http.Response, error) {
	if !r.Available() {
		return nil, nil
	}

	return r.resp, r.err
}

// HandleAction implements influx.Handler for the http system.
func (state *Result) HandleAction(
	ctx context.Context,
	action influx.Action,
) error {

	switch action := action.(type) {
	// a Response stored in the state tree copies the value of *Response
	// (disptached as an action against the store) when their IDs match.  While
	// this happens too late for parent components to respond to, this component
	// provides a cache of the response for later use.
	case *influx.Result:
		if !action.WillComplete(action) {
			return nil
		}

		resp, err := state.Complete(action)
		log.Println("resp:", resp, err)
	}

	return nil
}

// MarshalJSON implements json.Marshaler
func (r *Result) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
	// req := r.Request()

	// // if protected, just return nil
	// if !req.Unprotected {
	// }

	// var state struct {
	// 	ID    string
	// 	State interface{}
	// }
	// resp, err := r.Get()
	// switch {
	// case err != nil:
	// 	state.State = err.Error()
	// case resp != nil:
	// 	state.State = resp.Status
	// default:
	// 	state.State = "pending"
	// }

	// enc, err := json.Marshal(state)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "encode http state failed")
	// }

	// return enc, nil
}

var _ influx.Handler = &Result{}
var _ json.Marshaler = &Result{}
