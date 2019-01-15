package message

import (
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

func (ma *MockAdapter) MoveToDLQ(m Message) error {
	args := ma.Called(m)
	return args.Error(0)
}
