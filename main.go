package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fredw/igti-aws-lambda-payments/pkg/adapter"
	"github.com/fredw/igti-aws-lambda-payments/pkg/handler"
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
	a := adapter.NewAdapter(c)
	h := handler.NewHandler(l, c, a)

	lambda.Start(h.Handler)
}
