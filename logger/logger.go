package logger

import (
	"fmt"

	"github.com/fredw/igti-aws-lambda-payments/config"
	log "github.com/sirupsen/logrus"
)

// NewLogger create a new logger instance or panic when the consumed configuration value is invalid
func NewLogger(c *config.Config) *log.Logger {
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("error to parse log level %s", err))
	}

	l := log.New()
	l.Level = level
	l.Formatter = &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "timestamp",
			log.FieldKeyMsg:  "message",
		},
	}

	return l
}
