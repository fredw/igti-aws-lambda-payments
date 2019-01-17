package provider

import (
	"reflect"

	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/stretchr/testify/mock"
)

// MockProvider represents a mocked provider
type MockProvider struct {
	mock.Mock
}

// Process mocks the process of the message
func (mp *MockProvider) Process(m message.Message) error {
	args := mp.Called(m)
	return args.Error(0)
}

// MockProviderList is a mocked list of the provider
type MockProviderList struct {
	mock.Mock
}

// GetByMessage mocks the return of the provider by a message
func (mpl *MockProviderList) GetByMessage(m message.Message) Processor {
	args := mpl.Called(m)
	arg := args.Get(0)
	if reflect.ValueOf(arg).IsNil() {
		return nil
	}
	return arg.(Processor)
}

// GetNames mocks the return of provider names
func (mpl *MockProviderList) GetNames() []string {
	args := mpl.Called()
	return args.Get(0).([]string)
}
