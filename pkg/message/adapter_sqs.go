package message

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
)

type SQSManager interface {
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

type SQSAdapter struct {
	Config *config.Config
	SQS    SQSManager
}

// NewSQSAdapter creates a new SQS adapter
func NewSQSAdapter(c *config.Config, sqs SQSManager) *SQSAdapter {
	a := &SQSAdapter{
		Config: c,
		SQS:    sqs,
	}
	return a
}

// GetMessages returns messages from SQS
func (a *SQSAdapter) GetMessages() (Messages, error) {
	result, err := a.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &a.Config.SqsQueueURL,
		MaxNumberOfMessages: &a.Config.SqsMaxNumberOfMessages,
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to read messages from SQS")
	}

	messages := Messages{}
	for _, rm := range result.Messages {
		m := Message{Id: rm.ReceiptHandle}
		b := *rm.Body
		if err := json.Unmarshal([]byte(b), &m); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

// Delete message from SQS
func (a *SQSAdapter) Delete(id *string) error {
	_, err := a.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &a.Config.SqsQueueURL,
		ReceiptHandle: id,
	})

	if err != nil {
		return errors.Wrap(err, "failed to delete messages from SQS")
	}

	return nil
}

// MoveToDLQ moves the message directly to the DLQ
func (a *SQSAdapter) MoveToDLQ(m Message) error {
	body, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}

	// Send the message to the DLQ
	id := string(uuid.NewV4().String())
	_, err = a.SQS.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(string(body)),
		QueueUrl:               aws.String(a.Config.SqsDLQQueueURL),
		MessageGroupId:         &id,
		MessageDeduplicationId: &id,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create the on the DLQ")
	}

	// Delete the message from the main SQS
	if err := a.Delete(m.Id); err != nil {
		return errors.Wrap(err, "failed to delete the message from the main SQS")
	}

	return nil
}
