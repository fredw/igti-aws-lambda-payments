package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/handler"
	"github.com/fredw/igti-aws-lambda-payments/pkg/logger"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
)

func main() {
	c, err := config.Load()
	if err != nil {
		panic("cannot load config")
	}

	l := logger.NewLogger(c)
	l.Info("application started successfully")
	l.WithField("config", c).Info("loaded config")

	// Create a list with all available providers
	providers := provider.NewProviders(c)
	l.WithField("providers", providers.GetNames()).Info("providers list")

	// Create a new SQS adapter
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	adapter := message.NewSQSAdapter(c, sqs.New(sess))

	// Create a new handler to handle the Lambda invocation
	h := handler.NewHandler(l, providers, adapter)

	lambda.Start(h.Handler)
}
