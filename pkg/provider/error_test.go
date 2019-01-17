package provider_test

import (
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func TestNewCriticalError(t *testing.T) {
	err := provider.NewCriticalError("test")
	assert.IsType(t, &provider.CriticalError{}, err)
	assert.Equal(t, "test", err.Error())
}
