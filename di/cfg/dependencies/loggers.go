package dependencies

import (
	"context"
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/lib/logger"
	"github.com/sarulabs/dingo/v4"
	"github.com/sirupsen/logrus"
)

const (
	loggerLogrus = "logger_logrus"
)

func GetLoggers() []dingo.Def {
	return []dingo.Def{
		{
			Name: loggerLogrus,
			Build: func(cfg *config.App) (*logrus.Entry, error) {
				return logger.New(context.TODO()), nil
			},
		},
	}
}
