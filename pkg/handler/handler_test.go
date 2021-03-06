package handler_test

import (
	"context"
	"io/ioutil"
	"testing"

	perrors "github.com/fredw/igti-aws-lambda-payments/pkg/errors"
	"github.com/fredw/igti-aws-lambda-payments/pkg/handler"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	messageID := "message-id"
	messages := message.Messages{
		{
			Id:       &messageID,
			Provider: "Example",
		},
	}

	tests := []struct {
		name                      string
		adapterGetMessageResponse message.Messages
		adapterGetMessageError    error
		adapterDeleteError        error
		adapterMoveDLQError       error
		processError              error
		providerEmpty             bool
		wantResponse              handler.Response
		wantErr                   error
	}{
		{
			name: "messages processed successful with 2 success",
			adapterGetMessageResponse: message.Messages{
				{
					Id:       &messageID,
					Provider: "Example",
				},
				{
					Id:       &messageID,
					Provider: "Example",
				},
			},
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						ID:     &messageID,
						Status: handler.MessageStatusSuccess,
					},
					{
						ID:     &messageID,
						Status: handler.MessageStatusSuccess,
					},
				},
			},
		},
		{
			name:                   "failed to return messages",
			adapterGetMessageError: errors.New("test"),
			wantErr:                handler.ErrFailedReadMessages,
		},
		{
			name:                      "no messages",
			adapterGetMessageResponse: message.Messages{},
			wantResponse: handler.Response{
				Result: "No messages received",
			},
		},
		{
			name:                      "messages processed with error",
			adapterGetMessageResponse: messages,
			processError:              errors.New("test"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						ID:     &messageID,
						Status: handler.MessageStatusError,
						Error:  "failed to process the payment: test",
					},
				},
			},
		},
		{
			name:                      "messages processed with error by non existent provider",
			adapterGetMessageResponse: messages,
			providerEmpty:             true,
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						ID:     &messageID,
						Status: handler.MessageStatusError,
						Error:  "provider Example not available to process this message",
					},
				},
			},
		},
		{
			name:                      "messages processed with critical error",
			adapterGetMessageResponse: messages,
			processError:              perrors.NewCriticalError("test"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						ID:     &messageID,
						Status: handler.MessageStatusCritical,
						Error:  "test",
					},
				},
			},
		},
		{
			name:                      "messages processed with delete error",
			adapterGetMessageResponse: messages,
			adapterDeleteError:        perrors.NewCriticalError("failed to delete messages from SQS"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						ID:     &messageID,
						Status: handler.MessageStatusCritical,
						Error:  "failed to delete messages from SQS",
					},
				},
			},
		},
		{
			name:                      "messages processed with DLQ error",
			adapterGetMessageResponse: messages,
			processError:              perrors.NewCriticalError("test"),
			adapterMoveDLQError:       errors.New("test"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						ID:     &messageID,
						Status: handler.MessageStatusError,
						Error:  "problem to move the message to DLQ: test",
					},
				},
			},
		},
	}

	ctx := context.TODO()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			l := log.New()
			l.Out = ioutil.Discard

			providerMock := new(provider.MockProvider)
			providerMock.On("Process", mock.AnythingOfType("message.Message")).Return(tc.processError)

			providerReturn := providerMock
			if tc.providerEmpty {
				providerReturn = nil
			}

			providersMock := new(provider.MockProviderList)
			providersMock.On("GetByMessage", mock.AnythingOfType("message.Message")).Return(providerReturn)

			mockAdapter := new(message.MockAdapter)
			mockAdapter.On("GetMessages").Return(tc.adapterGetMessageResponse, tc.adapterGetMessageError)
			mockAdapter.On("Delete", mock.Anything).Return(tc.adapterDeleteError)
			mockAdapter.On("MoveToFailed", mock.Anything).Return(tc.adapterMoveDLQError)

			h := handler.NewHandler(l, providersMock, mockAdapter)
			resp, err := h.Handler(ctx, handler.Event{})

			assert.Equal(t, tc.wantResponse, resp)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
