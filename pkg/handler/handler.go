package handler

import (
	"context"
	"fmt"

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

type Event struct{}

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
	Id     *string `json:"id"`
	Status string  `json:"status"`
	Error  string  `json:"error"`
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

		mStatus := MessageStatusSuccess
		mError := ""
		if err != nil {
			mStatus = MessageStatusError
			switch err.(type) {
			case *provider.CriticalError:
				mStatus = MessageStatusCritical
			}
			mError = err.Error()
			h.log.WithError(err).WithField("message", m).Info("problem to process message")
		}

		mr := MessageResponse{
			Id:     m.Id,
			Status: mStatus,
			Error:  mError,
		}
		mResponses = append(mResponses, mr)

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
		return errors.New(fmt.Sprintf("provider %s not available to process this message", m.Provider))
	}

	// Try to process the message
	if err := p.Process(m); err != nil {
		// If it's a critical failure, move the message directly to the DLQ
		switch err.(type) {
		case *provider.CriticalError:
			errDLQ := h.adapter.MoveToDLQ(m)
			if errDLQ != nil {
				return errors.Wrap(err, "problem to move the message to DLQ")
			}
			return err
		}
		return errors.Wrap(err, "failed to process the payment")
	}

	// Successful: delete the message from SQS
	if err := h.adapter.Delete(m.Id); err != nil {
		return errors.Wrap(err, "failed to delete messages from SQS")
	}

	return nil
}
