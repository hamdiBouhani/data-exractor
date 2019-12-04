package common

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerContext struct{}

var LogCtx = loggerContext{}

func Log(ctx context.Context) logrus.FieldLogger {
	logger, ok := ctx.Value(LogCtx).(logrus.FieldLogger)
	if !ok {
		logger = logrus.New()
	}
	return logger
}
