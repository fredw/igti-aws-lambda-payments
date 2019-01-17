package message_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSQSAdapter_GetMessages(t *testing.T) {

	messageId := "123"

	tests := []struct {
		name                 string
		receiveMessageOutput *sqs.ReceiveMessageOutput
		receiveMessageError  error
		want                 message.Messages
		wantError            error
		wantErrorType        interface{}
	}{
		{
			name: "returned messages successfully",
			receiveMessageOutput: &sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						MessageId:     aws.String("123"),
						ReceiptHandle: aws.String("123"),
						Body:          aws.String(`{"provider":"test"}`),
					},
				},
			},
			want: message.Messages{
				message.Message{
					Id:       &messageId,
					Provider: "test",
					Order:    message.Order{},
				},
			},
		},
		{
			name:                "failed by SQS received messages",
			receiveMessageError: errors.New("test"),
			wantError:           errors.New("test"),
		},
		{
			name: "failed by unmarshal message body",
			receiveMessageOutput: &sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						MessageId:     aws.String("123"),
						ReceiptHandle: aws.String("123"),
						Body:          aws.String(`this is not a valid json body`),
					},
				},
			},
			wantErrorType: &json.SyntaxError{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSQS := new(message.MockSQS)
			mockSQS.On("ReceiveMessage", mock.AnythingOfType("*sqs.ReceiveMessageInput")).
				Return(tc.receiveMessageOutput, tc.receiveMessageError)

			sa := message.NewSQSAdapter(&config.Config{}, mockSQS)
			messages, err := sa.GetMessages()

			assert.Equal(t, tc.want, messages)
			if err != nil && tc.wantErrorType != nil {
				assert.IsType(t, tc.wantErrorType, err)
			}
		})
	}
}

func TestSQSAdapter_Delete(t *testing.T) {

	messageId := "123"

	tests := []struct {
		name        string
		messageId   *string
		deleteError error
		wantError   bool
	}{
		{
			name:      "message deleted successfully",
			messageId: &messageId,
		},
		{
			name:        "message deleted failed",
			messageId:   &messageId,
			deleteError: errors.New("test"),
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSQS := new(message.MockSQS)
			mockSQS.On("DeleteMessage", mock.AnythingOfType("*sqs.DeleteMessageInput")).
				Return(nil, tc.deleteError)

			sa := message.NewSQSAdapter(&config.Config{}, mockSQS)
			err := sa.Delete(tc.messageId)

			if tc.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSQSAdapter_MoveToFailed(t *testing.T) {

	messageId := "123"

	tests := []struct {
		name        string
		message     message.Message
		sendError   error
		deleteError error
		wantError   bool
	}{
		{
			name: "message moved successfully",
			message: message.Message{
				Id: &messageId,
			},
		},
		{
			name: "message moved failed due the send error",
			message: message.Message{
				Id: &messageId,
			},
			sendError: errors.New("test"),
			wantError: true,
		},
		{
			name: "message moved failed due the delete error",
			message: message.Message{
				Id: &messageId,
			},
			deleteError: errors.New("test"),
			wantError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSQS := new(message.MockSQS)
			mockSQS.On("SendMessage", mock.AnythingOfType("*sqs.SendMessageInput")).
				Return(nil, tc.sendError)
			mockSQS.On("DeleteMessage", mock.AnythingOfType("*sqs.DeleteMessageInput")).
				Return(nil, tc.deleteError)

			sa := message.NewSQSAdapter(&config.Config{}, mockSQS)
			err := sa.MoveToFailed(tc.message)

			if tc.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
