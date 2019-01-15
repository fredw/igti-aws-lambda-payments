package provider

import (
	"testing"

	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/stretchr/testify/assert"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
)

var providers = Providers{}

func init() {
	c, _ := config.Load()
	providers = NewProviders(c)
}

func TestProviders(t *testing.T) {
	assert.Equal(t, 1, len(providers))
}

func TestGetByMessage(t *testing.T) {
	tests := []struct {
		name    string
		message message.Message
		want    Processor
	}{
		{
			name: "it should return the example provider",
			message: message.Message{
				Provider: "Example",
			},
			want: providers[ExampleProvider],
		},
		{
			name: "it should't return a provider",
			message: message.Message{
				Provider: "UnkwownProvider",
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			provider := providers.GetByMessage(tc.message)
			assert.Equal(t, tc.want, provider)
		})
	}
}

func TestGetNames(t *testing.T) {
	want := []string{"Example"}
	assert.Equal(t, want, providers.GetNames())
}
