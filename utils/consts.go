package utils

import (
	"context"

	"spf-playlist/pkg/logger"
)

func GetLogger(ctx context.Context) logger.Logger {
	return ctx.Value("logger").(logger.Logger)
}
