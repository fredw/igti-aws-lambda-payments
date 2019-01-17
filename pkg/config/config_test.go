package config_test

import (
	"os"
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    *config.Config
		wantErr string
	}{
		{
			name: "environment variable are set correctly",
			env: map[string]string{
				"LOG_LEVEL":                    "INFO",
				"SQS_QUEUE_URL":                "http://sqs.host/",
				"SQS_DLQ_QUEUE_URL":            "http://sqs.dlq.host/",
				"SQS_MAX_NUMBER_OF_MESSAGES":   "1",
				"PROVIDER_EXAMPLE_REQUEST_URI": "http://provider.host/",
			},
			want: &config.Config{
				LogLevel:                  "INFO",
				SqsQueueURL:               "http://sqs.host/",
				SqsDLQQueueURL:            "http://sqs.dlq.host/",
				SqsMaxNumberOfMessages:    1,
				ProviderExampleRequestURI: "http://provider.host/",
			},
		},
		{
			name: "SQS_QUEUE_URL missing",
			env: map[string]string{
				"LOG_LEVEL":                    "debug",
				"SQS_DLQ_QUEUE_URL":            "http://sqs.dlq.host/",
				"SQS_MAX_NUMBER_OF_MESSAGES":   "1",
				"PROVIDER_EXAMPLE_REQUEST_URI": "http://provider.host/",
			},
			wantErr: "required key SQS_QUEUE_URL missing value",
		},
		{
			name: "SQS_DLQ_QUEUE_URL missing",
			env: map[string]string{
				"LOG_LEVEL":                    "debug",
				"SQS_QUEUE_URL":                "http://sqs.host/",
				"SQS_MAX_NUMBER_OF_MESSAGES":   "1",
				"PROVIDER_EXAMPLE_REQUEST_URI": "http://provider.host/",
			},
			wantErr: "required key SQS_DLQ_QUEUE_URL missing value",
		},
		{
			name: "PROVIDER_EXAMPLE_REQUEST_URI missing",
			env: map[string]string{
				"LOG_LEVEL":                  "debug",
				"SQS_QUEUE_URL":              "http://sqs.host/",
				"SQS_DLQ_QUEUE_URL":          "http://sqs.dlq.host/",
				"SQS_MAX_NUMBER_OF_MESSAGES": "1",
			},
			wantErr: "required key PROVIDER_EXAMPLE_REQUEST_URI missing value",
		},
		{
			name: "incorrect int var",
			env: map[string]string{
				"LOG_LEVEL":                    "INFO",
				"SQS_QUEUE_URL":                "http://sqs.host/",
				"SQS_DLQ_QUEUE_URL":            "http://sqs.dlq.host/",
				"SQS_MAX_NUMBER_OF_MESSAGES":   "test",
				"PROVIDER_EXAMPLE_REQUEST_URI": "http://provider.host/",
			},
			wantErr: "envconfig.Process: assigning SQS_MAX_NUMBER_OF_MESSAGES to SqsMaxNumberOfMessages: converting 'test' to type int64. details: strconv.ParseInt: parsing \"test\": invalid syntax",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			setEnv(tc.env)

			c, err := config.Load()

			assert.Equal(t, tc.want, c)
			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
			}
		})
	}
}

// setEnv sets environment variables
func setEnv(env map[string]string) {
	os.Clearenv()
	for key, val := range env {
		_ = os.Setenv(key, val)
	}
}
