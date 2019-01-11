package adapter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
)

type SQSAdapter struct {
	config *config.Config
	svc    *sqs.SQS
}

func NewAdapter(c *config.Config) *SQSAdapter {
	a := &SQSAdapter{
		config: c,
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	a.svc = sqs.New(sess)
	return a
}

// GetMessages returns messages from SQS
func (a *SQSAdapter) GetMessages() (*sqs.ReceiveMessageOutput, error) {
	result, err := a.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &a.config.SqsQueueURL,
		MaxNumberOfMessages: &a.config.SqsMaxNumberOfMessages,
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to read messages from SQS")
	}

	return result, nil
}

// Delete message from SQS
func (a *SQSAdapter) Delete(rh *string) error {
	_, err := a.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &a.config.SqsQueueURL,
		ReceiptHandle: rh,
	})

	if err != nil {
		return errors.Wrap(err, "failed to delete messages from SQS")
	}

	return nil
}
