package provider

import (
	"net/http"
	"time"

	"github.com/fredw/igti-aws-lambda-payments/pkg/client"
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/pkg/errors"
)

var (
	ErrFailProcessPayment       = errors.New("fail to process the payment")
	ErrFailedRequest            = errors.New("failed to do a request to the provider")
	ErrCriticalProviderInternal = NewCriticalError("payment can't be processed due a provider internal error")
)

// Example represents an example provider
type Example struct {
	Config *config.Config
	Client client.HttpCaller
}

// NewExampleProvider returns a new example provider
func NewExampleProvider(config *config.Config) Example {
	// Calculate the request timeout
	timeout := time.Duration(time.Duration(60) * time.Second)
	// Create a http client
	c := &http.Client{Timeout: timeout}

	p := Example{
		Config: config,
		Client: client.NewHttpClient(c),
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
		return ErrFailedRequest

	}
	if err = resp.Body.Close(); err != nil {
		return errors.Wrap(err, "error on close response body")
	}

	// Critical failure on provider, the message shouldn't be processed again, moving directly to the DLQ
	// For example, you can check for a specific error on message body. In this case we are checking for 500 Internal Server Error
	if resp.StatusCode == http.StatusInternalServerError {
		return ErrCriticalProviderInternal
	}

	// Payment failed on provider
	// For example, this provider consider a payment failure when the http status is different from 200 OK
	if resp.StatusCode != http.StatusOK {
		return ErrFailProcessPayment
	}

	return nil
}
