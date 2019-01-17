package message

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/mock"
)

// MockAdapter represents a mocked adapter
type MockAdapter struct {
	mock.Mock
}

// GetMessages mocks the return of the messages
func (ma *MockAdapter) GetMessages() (Messages, error) {
	args := ma.Called()
	return args.Get(0).(Messages), args.Error(1)
}

// GetMessages mocks the message deletion
func (ma *MockAdapter) Delete(id *string) error {
	args := ma.Called(id)
	return args.Error(0)
}

// MoveToFailed mocks the message being moved to failed
func (ma *MockAdapter) MoveToFailed(m Message) error {
	args := ma.Called(m)
	return args.Error(0)
}

// MockSQS represents a mocked SQS manager
type MockSQS struct {
	mock.Mock
}

// ReceiveMessage mocks the receive message
func (ms *MockSQS) ReceiveMessage(rmi *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := ms.Called(rmi)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

// DeleteMessage mocks the delete message
func (ms *MockSQS) DeleteMessage(dmi *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	args := ms.Called(dmi)
	return nil, args.Error(1)
}

// SendMessage mocks the send message
func (ms *MockSQS) SendMessage(smi *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := ms.Called(smi)
	return nil, args.Error(1)
}
