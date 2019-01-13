package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents common application parameters
type Config struct {
	LogLevel                  string `envconfig:"LOG_LEVEL" default:"INFO"`
	SqsQueueURL               string `envconfig:"SQS_QUEUE_URL" required:"true"`
	SqsDLQQueueURL            string `envconfig:"SQS_DLQ_QUEUE_URL" required:"true"`
	SqsMaxNumberOfMessages    int64  `envconfig:"SQS_MAX_NUMBER_OF_MESSAGES" default:"1"`
	ProviderExampleRequestURI string `envconfig:"PROVIDER_EXAMPLE_REQUEST_URI" required:"true"`
}

// Load loads the environment variables
func Load() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
