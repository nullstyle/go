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
) (influx.Request, error) {

	store, err := influx.FromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get store failed")
	}

	req := influx.NewRequest()

	stdreq, err := http.NewRequest("GET", url, nil)
	stdreq = stdreq.WithContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "make request failed")
	}

	store.Go(func() {
		resp, err := client.Raw.Do(stdreq)
		ctx = context.WithValue(ctx, req, result{resp, err})
		store.Dispatch(ctx, req)
	})

	return req, nil
}
