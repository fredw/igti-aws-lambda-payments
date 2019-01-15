package client

import (
	"errors"
	"net/http"
)

var (
	ErrMalformedBody = errors.New("malformed body")
)

// HttpCaller representation of the client call
type HttpCaller interface {
	Do(r *http.Request) (*http.Response, error)
}

// HttpClient is the wrapper of the http client with the purpose of having a mockable wrapper around it
type HttpClient struct {
	client HttpCaller
}

// NewHttpClient create a new instance of HttpCaller
func NewHttpClient(c HttpCaller) *HttpClient {
	return &HttpClient{c}
}

// Do Perform a request returning the response this wrapper was creating for testing purposes and to have
// an easier way to switch from the native client in case needed
func (c HttpClient) Do(r *http.Request) (*http.Response, error) {
	res, err := c.client.Do(r)

	if err != nil {
		return nil, err
	}

	if res.Body == nil {
		return nil, ErrMalformedBody
	}

	return res, nil
}
