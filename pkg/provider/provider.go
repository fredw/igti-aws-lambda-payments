package provider

import (
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
)

// Available providers
const (
	providerExample = "Example"
)

// Providers represents a list of providers
type Providers map[string]Processor

// Processor represents a provider that can process a message
type Processor interface {
	Process(m message.Message) error
}

// NewProviders create a list of all available providers
func NewProviders(config *config.Config) Providers {
	providers := Providers{
		providerExample: NewExampleProvider(config),
	}
	return providers
}

// GetByMessage return a provider checking the the provider string on the message
func (ps Providers) GetByMessage(m message.Message) Processor {
	return ps[m.Provider]
}
