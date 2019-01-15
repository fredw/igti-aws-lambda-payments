package provider

import (
	"reflect"

	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/stretchr/testify/mock"
)

type MockProvider struct {
	mock.Mock
}

func (mp *MockProvider) Process(m message.Message) error {
	args := mp.Called(m)
	return args.Error(0)
}

type MockProviderList struct {
	mock.Mock
}

func (mpl *MockProviderList) GetByMessage(m message.Message) Processor {
	args := mpl.Called(m)
	arg := args.Get(0)
	if reflect.ValueOf(arg).IsNil() {
		return nil
	}
	return arg.(Processor)
}

func (mpl *MockProviderList) GetNames() []string {
	args := mpl.Called()
	return args.Get(0).([]string)
}
