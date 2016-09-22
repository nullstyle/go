package http

import (
	"net/http"

	"github.com/nullstyle/go/influx"

	// HACK(scott): the import below works around what appears to be a glide bug,
	// resulting in dependencies that are only referenced in tests no being
	// detected properly.  It's less than ideal, but this hope is that this can be
	// removed before it becomes a liability, when glide behaves as expected.
	_ "github.com/stellar/go/support/http/httptest"
)

// DefaultClient is an influx service that powers the http system.  When
// components call to the service, the http request is made in the background
// then dispatched as an influx.HttpResponse action.
var DefaultClient = &Client{
	Raw: http.DefaultClient,
}

// Client implements an influx-aware http client.  Components can call into a clietnt
type Client struct {
	Raw    HTTP
	nextID int
}

// HTTP reperesents a type that can make http requests.
type HTTP interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request is an influx component that represents an http request.
type Result struct {

	// Request represents the request that originated this response.
	influx.Request

	result *influx.Result
}

var _ HTTP = http.DefaultClient
