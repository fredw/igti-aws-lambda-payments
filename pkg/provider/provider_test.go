package provider_test

import (
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/stretchr/testify/assert"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
)

var providers = provider.Providers{}

func init() {
	c, _ := config.Load()
	providers = provider.NewProviders(c)
}

func TestProviders(t *testing.T) {
	assert.Equal(t, 1, len(providers))
}

func TestGetByMessage(t *testing.T) {
	tests := []struct {
		name    string
		message message.Message
		want    provider.Processor
	}{
		{
			name: "it should return the example providerExample",
			message: message.Message{
				Provider: "Example",
			},
			want: providers[provider.ExampleProvider],
		},
		{
			name: "it should't return a providerExample",
			message: message.Message{
				Provider: "UnkwownProvider",
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := providers.GetByMessage(tc.message)
			assert.Equal(t, tc.want, p)
		})
	}
}

func TestGetNames(t *testing.T) {
	want := []string{"Example"}
	assert.Equal(t, want, providers.GetNames())
}
