package provider

import (
	"net/http"
	"time"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/pkg/errors"
)

// Example represents an example provider
type Example struct {
	Config *config.Config
	Client *http.Client
}

// NewExampleProvider returns a new example provider
func NewExampleProvider(config *config.Config) Example {
	// Calculate the request timeout
	timeout := time.Duration(time.Duration(60) * time.Second)
	// Create a http client
	client := &http.Client{Timeout: timeout}

	p := Example{
		Config: config,
		Client: client,
	}

	return p
}

func (p Example) Process(m message.Message) error {
	// Create a request to the provider
	req, err := http.NewRequest(http.MethodPost, p.Config.ProviderExampleRequestURI, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create a request")
	}

	// Do the request
	resp, err := p.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed do a request to the provider")

	}
	if err = resp.Body.Close(); err != nil {
		return errors.Wrap(err, "error on close response body")
	}

	// Critical failure on provider, the message shouldn't be processed again, moving directly to the DLQ
	// For example, you can check for a specific error on message body. In this case we are checking for 500 Internal Server Error
	if resp.StatusCode == http.StatusInternalServerError {
		return NewCriticalError("payment can't be processed due a provider internal error")
	}

	// Payment failed on provider
	// For example, this provider consider a payment failure when the http status is different from 200 OK
	if resp.StatusCode != http.StatusOK {
		return errors.New("fail to process the payment")
	}

	return nil
}