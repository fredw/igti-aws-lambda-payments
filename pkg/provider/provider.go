package provider

import (
	"net/http"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
)

// Available providers
const (
	ExampleProvider = "Example"
)

// HTTPClient representation of the client call
type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

// Providers represents a list of providers
type Providers map[string]Processor

// Processor represents a provider that can process a message
type Processor interface {
	Process(m message.Message) error
}

// NewProviders create a list of all available providers
func NewProviders(config *config.Config) Providers {
	providers := Providers{
		ExampleProvider: NewExampleProvider(config),
	}
	return providers
}

// GetByMessage returns a provider checking the the provider string on the message
func (providers Providers) GetByMessage(m message.Message) Processor {
	return providers[m.Provider]
}

// GetNames returns the providers names
func (providers Providers) GetNames() []string {
	var names []string
	for k := range providers {
		names = append(names, k)
	}
	return names
}
