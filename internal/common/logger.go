package common

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	RequestIdKey = "requestid"
)

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := createLogger(ctx)
	return logger
}

func GetRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, ok := ctx.Value(RequestIdKey).(string)
	if !ok {
		return ""
	}
	return id
}

func createLogger(ctx context.Context) *logrus.Entry {
	log := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			DisableQuote: true,
			ForceColors:  true,
		},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ReportCaller: false,
	}
	requestId := GetRequestId(ctx)
	if requestId == "" {
		return log.WithContext(ctx)
	}
	return log.WithContext(ctx).WithField(RequestIdKey, GetRequestId(ctx))
}
