package http

import (
	"net/http"
	"sync"
)

// result represents the protected internal result state of an http request
type result struct {
	sync.Mutex
	*http.Response
	Err error
}

func (result *result) get() (*http.Response, error) {
	result.Lock()
	response, err := result.Response, result.Err
	result.Unlock()
	return response, err
}

func (result *result) finish(resp *http.Response, err error) {
	result.Lock()
	result.Response = resp
	result.Err = err
	result.Unlock()
}
