package handler

import (
	"context"
	"fmt"

	perrors "github.com/fredw/igti-aws-lambda-payments/pkg/errors"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Errors
var (
	ErrFailedReadMessages = errors.New("failed to read messages from SQS")
)

// Message statuses
var (
	MessageStatusSuccess  = "success"
	MessageStatusError    = "error"
	MessageStatusCritical = "critical"
)

// Event represents the Lambda event
type Event struct{}

// Handler represents the handler
type Handler struct {
	log       *log.Logger
	providers provider.ProcessorList
	adapter   message.Adapter
}

// Response represents the lambda response
type Response struct {
	Result   string            `json:"result"`
	Messages []MessageResponse `json:"messages"`
}

// MessageResponse represents the message response
type MessageResponse struct {
	ID     *string `json:"id"`
	Status string  `json:"status"`
	Error  string  `json:"error,omitempty"`
}

// NewHandler creates a new handler struct
func NewHandler(l *log.Logger, p provider.ProcessorList, a message.Adapter) *Handler {
	h := &Handler{
		log:       l,
		providers: p,
		adapter:   a,
	}
	return h
}

// Handler handles the lambda invoke
func (h *Handler) Handler(ctx context.Context, event Event) (Response, error) {
	messages, err := h.adapter.GetMessages()

	if err != nil {
		return Response{}, ErrFailedReadMessages
	}
	if len(messages) == 0 {
		return Response{Result: "No messages received"}, nil
	}

	// Process all messages
	var mResponses []MessageResponse
	for _, m := range messages {
		err := h.processMessage(m)
		mResponses = append(mResponses, h.getMessageResponse(m, err))

		h.log.WithFields(log.Fields{"message": m}).Info("message processed successfully")
	}

	return Response{
		Result:   "Messages processed",
		Messages: mResponses,
	}, nil
}

// processMessage process a message calling the provider logic and handle the message through the SQS
func (h *Handler) processMessage(m message.Message) error {
	// Get the provider and process the message using the own provider logic
	p := h.providers.GetByMessage(m)
	if p == nil {
		return fmt.Errorf("provider %s not available to process this message", m.Provider)
	}

	// Try to process the message
	if err := p.Process(m); err != nil {
		return h.processErrorMessage(m, err)
	}

	// After successful process, try to delete the message from SQS
	if err := h.adapter.Delete(m.Id); err != nil {
		return h.processErrorMessage(m, err)
	}

	return nil
}

// processErrorMessage process a message with an error
func (h *Handler) processErrorMessage(m message.Message, err error) error {
	switch err.(type) {
	case *perrors.CriticalError:
		// If it's a critical failure, move the message directly to the failed list
		errM := h.adapter.MoveToFailed(m)
		if errM != nil {
			return errors.Wrap(err, "problem to move the message to DLQ")
		}
		return err
	}
	return errors.Wrap(err, "failed to process the payment")
}

// getMessageResponse returns a message response
func (h *Handler) getMessageResponse(m message.Message, err error) MessageResponse {
	if err != nil {
		mStatus := MessageStatusError
		switch err.(type) {
		case *perrors.CriticalError:
			mStatus = MessageStatusCritical
		}

		h.log.WithError(err).WithField("message", m).Info("problem to process message")

		return MessageResponse{
			ID:     m.Id,
			Status: mStatus,
			Error:  err.Error(),
		}
	}

	return MessageResponse{
		ID:     m.Id,
		Status: MessageStatusSuccess,
	}
}
