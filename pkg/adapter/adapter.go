package adapter

import "github.com/aws/aws-sdk-go/service/sqs"

type MessageAdapter interface {
	GetMessages() (*sqs.ReceiveMessageOutput, error)
	Delete(rh *string) error
}
