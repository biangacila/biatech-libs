package utils

import (
	"log/slog"
	"os"
)

func NewLoggerSlogInfo(err error) {
	NewLoggerSlog().Info(err.Error())
}
func NewLoggerSlogWarn(err error) {
	NewLoggerSlog().Warn(err.Error())
}
func NewLoggerSlogError(err error) {
	NewLoggerSlog().Error(err.Error())
}
func NewLoggerSlog() *slog.Logger {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(logHandler)
	return logger
}
