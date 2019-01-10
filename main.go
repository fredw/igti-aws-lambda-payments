package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"

	"github.com/fredw/igti-aws-lambda-payments/config"
	"github.com/fredw/igti-aws-lambda-payments/logger"
)

type Event struct {
}

var c *config.Config
var l *log.Logger

func init() {
	var err error
	c, err = config.Load()
	if err != nil {
		panic("cannot load config")
	}

	l = logger.NewLogger(c)
	l.Info("application started successfully")
}

func Handler(ctx context.Context, event Event) (string, error) {
	l.WithField("config", c).Info("loaded config")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &c.SqsQueueURL,
		MaxNumberOfMessages: &c.SqsMaxNumberOfMessages,
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		return "", errors.New("failed to read messages from SQS")
	}

	if len(result.Messages) == 0 {
		return "no messages received", nil
	}

	var messagesProcessed int64
	var messagesFailed int64

	// Calculate the request timeout
	timeout := time.Duration(time.Duration(c.ProviderRequestTimeout) * time.Second)
	client := &http.Client{Timeout: timeout}

	for _, m := range result.Messages {
		go func() {
			req, err := http.NewRequest(http.MethodPost, c.ProviderRequestURI, nil)
			if err != nil {
				l.WithField("message", m).Error("failed create a request")
			}

			resp, err := client.Do(req)
			if err != nil {
				l.WithField("message", m).Error("failed do a request to the provider")
			}
			defer resp.Body.Close()

			// Payment messagesProcessed successfully
			if resp.StatusCode != http.StatusOK {
				l.WithField("message", m).Info("fail to process the payment")
				messagesFailed = messagesFailed + 1
				return
			}

			d, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &c.SqsQueueURL,
				ReceiptHandle: m.ReceiptHandle,
			})

			if err != nil {
				l.WithField("message", m).Error("failed to delete messages from SQS")
			}

			messagesProcessed = messagesProcessed + 1
			l.WithField("output", d).Info("message deleted successfully")
		}()
	}

	return fmt.Sprintf("successfully: %d failed: %d", messagesProcessed, messagesFailed), nil
}

func main() {
	lambda.Start(Handler)
}
