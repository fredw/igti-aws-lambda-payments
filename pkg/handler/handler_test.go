package handler_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/handler"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHandler(t *testing.T) {
	messageId := "message-id"
	messages := message.Messages{
		{
			Id:       &messageId,
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
					Id:       &messageId,
					Provider: "Example",
				},
				{
					Id:       &messageId,
					Provider: "Example",
				},
			},
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						Id:     &messageId,
						Status: handler.MessageStatusSuccess,
					},
					{
						Id:     &messageId,
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
						Id:     &messageId,
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
						Id:     &messageId,
						Status: handler.MessageStatusError,
						Error:  "provider Example not available to process this message",
					},
				},
			},
		},
		{
			name:                      "messages processed with critical error",
			adapterGetMessageResponse: messages,
			processError:              provider.NewCriticalError("test"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						Id:     &messageId,
						Status: handler.MessageStatusCritical,
						Error:  "test",
					},
				},
			},
		},
		{
			name:                      "messages processed with delete error",
			adapterGetMessageResponse: messages,
			adapterDeleteError:        errors.New("test"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						Id:     &messageId,
						Status: handler.MessageStatusError,
						Error:  "failed to delete messages from SQS: test",
					},
				},
			},
		},
		{
			name:                      "messages processed with DLQ error",
			adapterGetMessageResponse: messages,
			processError:              provider.NewCriticalError("test"),
			adapterMoveDLQError:       errors.New("test"),
			wantResponse: handler.Response{
				Result: "Messages processed",
				Messages: []handler.MessageResponse{
					{
						Id:     &messageId,
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
			mockAdapter.On("MoveToDLQ", mock.Anything).Return(tc.adapterMoveDLQError)

			h := handler.NewHandler(l, providersMock, mockAdapter)
			resp, err := h.Handler(ctx, handler.Event{})

			assert.Equal(t, tc.wantResponse, resp)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

//func TestHandler_Handler(t *testing.T) {
//	type fields struct {
//		log       *log.Logger
//		providers provider.Providers
//		adapter   message.Adapter
//	}
//	type args struct {
//		ctx   context.Context
//		event Event
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    Response
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &Handler{
//				log:       tt.fields.log,
//				providers: tt.fields.providers,
//				adapter:   tt.fields.adapter,
//			}
//			got, err := h.Handler(tt.args.ctx, tt.args.event)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Handler.Handler() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Handler.Handler() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestHandler_processMessage(t *testing.T) {
//	type fields struct {
//		log       *log.Logger
//		providers provider.Providers
//		adapter   message.Adapter
//	}
//	type args struct {
//		m message.Message
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &Handler{
//				log:       tt.fields.log,
//				providers: tt.fields.providers,
//				adapter:   tt.fields.adapter,
//			}
//			if err := h.processMessage(tt.args.m); (err != nil) != tt.wantErr {
//				t.Errorf("Handler.processMessage() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
