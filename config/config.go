package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents common application parameters
type Config struct {
	LogLevel               string `envconfig:"LOG_LEVEL" default:"INFO"`
	ProviderRequestTimeout int32  `envconfig:"PROVIDER_REQUEST_TIMEOUT" default:"30"`
	ProviderRequestURI     string `envconfig:"PROVIDER_REQUEST_URI" required:"true"`
	SqsQueueURL            string `envconfig:"SQS_QUEUE_URL" required:"true"`
	SqsMaxNumberOfMessages int64  `envconfig:"SQS_MAX_NUMBER_OF_MESSAGES" default:"1"`
}

// Load loads the env
func Load() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
