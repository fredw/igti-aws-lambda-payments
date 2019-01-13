package handler

import (
	"context"
	"fmt"

	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Event struct{}

type Handler struct {
	log       *log.Logger
	providers provider.Providers
	adapter   message.Adapter
}

// NewHandler creates a new handler struct
func NewHandler(l *log.Logger, p provider.Providers, a message.Adapter) *Handler {
	h := &Handler{
		log:       l,
		providers: p,
		adapter:   a,
	}
	return h
}

// Response represents the lambda response
type Response struct {
	Message    string `json:"message"`
	Successful int64  `json:"successful"`
	Failed     int64  `json:"failed"`
}

// Handler handles the lambda invoke
func (h *Handler) Handler(ctx context.Context, event Event) (Response, error) {
	messages, err := h.adapter.GetMessages()

	if err != nil {
		return Response{}, errors.New("failed to read messages from SQS")
	}
	if len(messages) == 0 {
		return Response{Message: "No messages received"}, nil
	}

	var successful int64
	var failed int64

	// Process all messages
	for _, m := range messages {
		if err := h.processMessage(m); err != nil {
			failed = failed + 1
			h.log.WithError(err).WithField("message", m).Info("problem to process message")
			continue
		}

		successful = successful + 1
		h.log.WithFields(log.Fields{"message": m}).Info("message processed successfully")
	}

	return Response{
		Message:    "Messages processed",
		Successful: successful,
		Failed:     failed,
	}, nil
}

// processMessage process a message calling the provider logic and handle the message through the SQS
func (h *Handler) processMessage(m message.Message) error {
	// Get the provider and process the message using the own provider logic
	p := h.providers.GetByMessage(m)
	if p == nil {
		return errors.New(fmt.Sprintf("provider %s not available to process this message", m.Provider))
	}

	// Try to process the message
	if err := p.Process(m); err != nil {
		// If it's a critical failure, move the message directly to the DLQ
		switch err.(type) {
		case *provider.CriticalError:
			err = h.adapter.MoveToDLQ(m)
		}
		return errors.Wrap(err, "failed to process the payment")
	}

	// Successful: delete the message from SQS
	if err := h.adapter.Delete(m.Id); err != nil {
		return errors.Wrap(err, "failed to delete messages from SQS")
	}

	return nil
}
