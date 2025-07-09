package common

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}

// Context'ten request-id ve trace-id alıp log entry'ye ekler
func WithContext(ctx context.Context) *logrus.Entry {
	reqID := ctx.Value("request-id")
	traceID := ctx.Value("trace-id")
	return Logger.WithFields(logrus.Fields{
		"request_id": reqID,
		"trace_id":   traceID,
	})
}
