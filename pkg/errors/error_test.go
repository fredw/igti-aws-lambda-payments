package errors_test

import (
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewCriticalError(t *testing.T) {
	err := errors.NewCriticalError("test")
	assert.IsType(t, &errors.CriticalError{}, err)
	assert.Equal(t, "test", err.Error())
}
