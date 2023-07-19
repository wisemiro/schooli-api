package logger

import (
	"os"

	"golang.org/x/exp/slog"
)

func GetLogger() *slog.Logger {
	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       nil,
		ReplaceAttr: nil,
	}))

	slog.SetDefault(l)
	return l
}
