package logger_test

import (
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	c := &config.Config{
		LogLevel: "INFO",
	}
	l := logger.NewLogger(c)

	assert.NotNil(t, l)
}
