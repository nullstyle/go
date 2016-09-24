package http

import "net/http"

// result represents the protected internal result state of an http request
type result struct {
	*http.Response
	Err error
}
