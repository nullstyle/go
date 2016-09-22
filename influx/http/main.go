package http

import (
	stdhttp "net/http"

	"github.com/nullstyle/go/influx"

	// HACK(scott): the import below works around what appears to be a glide bug,
	// resulting in dependencies that are only referenced in tests no being
	// detected properly.
	_ "github.com/stellar/go/support/http/httptest"
)

// DefaultClient is an influx service that powers the http system.  When
// components call to the service, the http request is made in the background
// then dispatched as an influx.HttpResponse action.
var DefaultClient = &Client{
	Raw: stdhttp.DefaultClient,
}

// Client implements an influx-aware http client.  Components can call into a clietnt
type Client struct {
	Raw    HTTP
	nextID int
}

// HTTP reperesents a type that can make http requests.
type HTTP interface {
	Do(req *stdhttp.Request) (*stdhttp.Response, error)
}

// Response is an influx component that represents an http request.
type Response struct {

	// Request represents the request that originated this response.
	Request influx.Request

	// Done is true when the result has been delivered
	Done bool

	// Resp is the response to the request.  When non-nil, the request is complete
	// and did not error.
	Resp *stdhttp.Response

	// Err is any error that occurred in when initiated the request. When non-nil,
	// an error occurred at the http client when making the request.
	Err error
}

var _ HTTP = stdhttp.DefaultClient
