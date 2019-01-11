package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fredw/igti-aws-lambda-payments/pkg/adapter"
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Event struct{}

type Handler struct {
	log     *log.Logger
	config  *config.Config
	adapter adapter.MessageAdapter
}

func NewHandler(l *log.Logger, c *config.Config, a adapter.MessageAdapter) *Handler {
	h := &Handler{
		log:     l,
		config:  c,
		adapter: a,
	}
	return h
}

type Response struct {
	Message    string `json:"message"`
	Successful int64  `json:"successful"`
	Failed     int64  `json:"failed"`
}

func (h *Handler) Handler(ctx context.Context, event Event) (Response, error) {
	h.log.WithField("config", h.config).Info("loaded config")

	result, err := h.adapter.GetMessages()

	if err != nil {
		return Response{}, errors.New("failed to read messages from SQS")
	}

	if len(result.Messages) == 0 {
		return Response{Message: "No messages received"}, nil
	}

	var successful int64
	var failed int64

	// Calculate the request timeout
	timeout := time.Duration(time.Duration(h.config.ProviderRequestTimeout) * time.Second)
	// Create a http client
	client := &http.Client{Timeout: timeout}

	// Process all returned messages
	for _, m := range result.Messages {
		err := h.processMessage(client, m)

		if err != nil {
			failed = failed + 1
			h.log.WithError(err).WithField("message", m)
			continue
		}

		successful = successful + 1
		h.log.WithFields(log.Fields{"message": m}).Info("message deleted successfully")
	}

	return Response{
		Message:    "Messages processed",
		Successful: successful,
		Failed:     failed,
	}, nil
}

func (h *Handler) processMessage(client *http.Client, m *sqs.Message) error {
	// Create a request to the provider
	req, err := http.NewRequest(http.MethodPost, h.config.ProviderRequestURI, nil)
	if err != nil {
		return errors.Wrap(err, "failed create a request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed do a request to the provider")

	}
	if err = resp.Body.Close(); err != nil {
		return errors.Wrap(err, "error on close response body")
	}

	// Payment failed on provider
	if resp.StatusCode != http.StatusOK {
		return errors.New("fail to process the payment")
	}

	err = h.adapter.Delete(m.ReceiptHandle)

	if err != nil {
		return errors.New("failed to delete messages from SQS")
	}

	return nil
}
