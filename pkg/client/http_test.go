package client_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/client"
	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name        string
		requestBody *bytes.Reader
		response    *http.Response
		responseErr error
		wantErr     error
	}{
		{
			name:        "successful request",
			requestBody: bytes.NewReader([]byte("boo")),
			response: &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`)),
			},
			responseErr: nil,
			wantErr:     nil,
		},
		{
			name:        "failed by response error",
			requestBody: bytes.NewReader([]byte("boo")),
			response:    nil,
			responseErr: errors.New("test"),
			wantErr:     errors.New("test"),
		},
		{
			name:        "failed by empty body",
			requestBody: bytes.NewReader([]byte("boo")),
			response: &http.Response{
				Body: nil,
			},
			wantErr: client.ErrMalformedBody,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := new(client.MockHTTPClient)
			req, _ := http.NewRequest(http.MethodGet, "/", tc.requestBody)
			mock.On("Do", req).Return(tc.response, tc.responseErr)

			c := client.NewHttpClient(mock)

			res, err := c.Do(req)
			if err == nil {
				assert.Equal(t, tc.response, res)
			}
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
