package http

import (
	"context"
	http "net/http"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// Get makes a get request
func (client *Client) Get(
	ctx context.Context,
	url string,
) (*Request, error) {

	store, err := influx.FromContext(ctx)
	if err != nil {
		return &Request{}, errors.Wrap(err, "get store failed")
	}

	var result result

	req := &Request{
		Request: influx.Request{
			ID: client.nextID,
		},
		result: &result,
	}
	client.nextID++

	stdreq, err := http.NewRequest("GET", url, nil)
	stdreq = stdreq.WithContext(ctx)
	if err != nil {
		return &Request{}, errors.Wrap(err, "make request failed")
	}

	store.Go(func() {
		resp, err := client.Raw.Do(stdreq)
		result.finish(resp, err)

		store.Dispatch(req)
	})

	return req, nil
}
