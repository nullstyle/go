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
) (Response, error) {

	store, err := influx.FromContext(ctx)
	if err != nil {
		return Response{}, errors.Wrap(err, "get store failed")
	}

	var result Response
	result.Request.ID = client.nextID
	client.nextID++

	req, err := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return Response{}, errors.Wrap(err, "make request failed")
	}

	store.Go(func() {
		resp, err := client.Raw.Do(req)
		if err != nil {
			result.Err = err
		} else {
			result.Resp = resp
		}

		store.Dispatch(&result)
	})

	return result, nil
}
