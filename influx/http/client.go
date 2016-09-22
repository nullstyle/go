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
) (Request, error) {

	store, err := influx.FromContext(ctx)
	if err != nil {
		return Request{}, errors.Wrap(err, "get store failed")
	}

	var result Request
	result.Request.ID = client.nextID
	client.nextID++

	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return Request{}, errors.Wrap(err, "make request failed")
	}

	store.Go(func() {
		resp, err := client.Raw.Do(req)
		result.finish(resp, err)

		store.Dispatch(&result)
	})

	return result, nil
}
