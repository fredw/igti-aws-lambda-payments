package message

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/mock"
)

type MockAdapter struct {
	mock.Mock
}

func (ma *MockAdapter) GetMessages() (Messages, error) {
	args := ma.Called()
	return args.Get(0).(Messages), args.Error(1)
}

func (ma *MockAdapter) Delete(id *string) error {
	args := ma.Called(id)
	return args.Error(0)
}

func (ma *MockAdapter) MoveToFailed(m Message) error {
	args := ma.Called(m)
	return args.Error(0)
}

type MockSQS struct {
	mock.Mock
}

func (ms *MockSQS) ReceiveMessage(rmi *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := ms.Called(rmi)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

func (ms *MockSQS) DeleteMessage(dmi *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	args := ms.Called(dmi)
	return nil, args.Error(1)
}

func (ms *MockSQS) SendMessage(smi *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	args := ms.Called(smi)
	return nil, args.Error(1)
}
