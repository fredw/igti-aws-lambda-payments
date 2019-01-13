package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fredw/igti-aws-lambda-payments/pkg/handler"
	"github.com/fredw/igti-aws-lambda-payments/pkg/message"
	"github.com/fredw/igti-aws-lambda-payments/pkg/provider"
	log "github.com/sirupsen/logrus"

	"github.com/fredw/igti-aws-lambda-payments/pkg/config"
	"github.com/fredw/igti-aws-lambda-payments/pkg/logger"
)

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

func main() {
	l.WithField("config", c).Info("loaded config")

	// Create a list with all available providers
	providers := provider.NewProviders(c)
	// Create a new adapter for SQS
	sqs := message.NewSQSAdapter(c)
	// Create a new handler to handle the Lambda invocation
	h := handler.NewHandler(l, providers, sqs)

	lambda.Start(h.Handler)
}
