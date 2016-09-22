package http

import (
	"context"
	"net/http"

	"github.com/nullstyle/go/influx"
	"github.com/pkg/errors"
)

// Get makes a get request
func (client *Client) Get(
	ctx context.Context,
	url string,
) (Result, error) {

	store, err := influx.FromContext(ctx)
	if err != nil {
		return Result{}, errors.Wrap(err, "get store failed")
	}

	stdreq, err := http.NewRequest("GET", url, nil)
	stdreq = stdreq.WithContext(ctx)
	if err != nil {
		return Result{}, errors.Wrap(err, "make request failed")
	}

	req, res := influx.NewRequest()

	store.Go(func() {
		resp, err := client.Raw.Do(stdreq)
		res.Finish(resp, err)

		store.Dispatch(influx.Done{req})
	})

	return Result{
		result: &res,
	}, nil
}
