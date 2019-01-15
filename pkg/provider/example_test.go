package provider

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/client"
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/stretchr/testify/assert"
)

var provider Example

var providerURI = "http://provider.host/"

func init() {
	c := &config.Config{
		ProviderExampleRequestURI: providerURI,
	}
	provider = NewExampleProvider(c)
}

func TestNewExampleProvider(t *testing.T) {
	assert.IsType(t, Example{}, provider)
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name          string
		message       message.Message
		response      *http.Response
		responseError error
		want          error
	}{
		//{
		//	name:    "failed by message without provider",
		//	message: message.Message{},
		//	response: &http.Response{
		//		StatusCode: http.StatusOK,
		//		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
		//	},
		//	want: ErrFailProcessPayment,
		//},
		//{
		//	name: "failed by message with an non existent provider",
		//	message: message.Message{
		//		Provider: "xyz",
		//	},
		//	response: &http.Response{
		//		StatusCode: http.StatusOK,
		//		Body:       ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
		//	},
		//	want: ErrFailProcessPayment,
		//},
		{
			name: "success due a 200 OK from provider",
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
			responseError: ErrFailedRequest,
			want:          ErrFailedRequest,
		},
		{
			name: "failed due a 400 Bad Request from provider",
			message: message.Message{
				Provider: "Example",
			},
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			},
			want: ErrFailProcessPayment,
		},
		{
			name: "failed due a 500 Internal Server Error from provider",
			message: message.Message{
				Provider: "Example",
			},
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			},
			want: ErrCriticalProviderInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mocked http client
			mock := new(client.MockHTTPClient)
			req, _ := http.NewRequest(http.MethodPost, providerURI, nil)
			mock.On("Do", req).Return(tc.response, tc.responseError)
			c := client.NewHttpClient(mock)

			// Overwrite the http client on provider
			provider.Client = c

			err := provider.Process(tc.message)
			assert.Equal(t, tc.want, err)
		})
	}
}
