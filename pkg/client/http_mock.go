package client

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockHTTPClient represents a mocked http client
type MockHTTPClient struct {
	mock.Mock
}

// Do does a http request
func (m *MockHTTPClient) Do(r *http.Request) (*http.Response, error) {
	args := m.Called(r)
	return args.Get(0).(*http.Response), args.Error(1)
}
