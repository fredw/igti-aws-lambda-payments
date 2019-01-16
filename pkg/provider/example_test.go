package provider_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/client"
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/stretchr/testify/assert"
)

var providerExample provider.Example

var providerURI = "http://providerExample.host/"

func init() {
	c := &config.Config{
		ProviderExampleRequestURI: providerURI,
	}
	providerExample = provider.NewExampleProvider(c)
}

func TestNewExampleProvider(t *testing.T) {
	assert.IsType(t, provider.Example{}, providerExample)
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name          string
		message       message.Message
		response      *http.Response
		responseError error
		want          error
	}{
		{
			name: "success due a 200 OK from providerExample",
			message: message.Message{
				Provider: "Example",
			},
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"result":"authorized"}`)),
			},
			want: nil,
		},
		{
			name: "failed due a response error",
			message: message.Message{
				Provider: "Example",
			},
			response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"error":"test"}`)),
			},
			responseError: provider.ErrFailedRequest,
			want:          provider.ErrFailedRequest,
		},
		{
			name: "failed due a 400 Bad Request from providerExample",
			message: message.Message{
				Provider: "Example",
			},
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			},
			want: provider.ErrFailProcessPayment,
		},
		{
			name: "failed due a 500 Internal Server Error from providerExample",
			message: message.Message{
				Provider: "Example",
			},
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
			want: provider.ErrCriticalProviderInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mocked http client
			mock := new(client.MockHTTPClient)
			req, _ := http.NewRequest(http.MethodPost, providerURI, nil)
			mock.On("Do", req).Return(tc.response, tc.responseError)
			c := client.NewHttpClient(mock)

			// Overwrite the http client on providerExample
			providerExample.Client = c

			err := providerExample.Process(tc.message)
			assert.Equal(t, tc.want, err)
		})
	}
}
