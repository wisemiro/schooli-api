package worker

import (
	"golang.org/x/exp/slog"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger(l *slog.Logger) *Logger {
	return &Logger{l: l}
}

func (logger *Logger) Print(args ...any) {
	logger.l.Debug("", args...)
}

func (logger *Logger) Debug(args ...any) {
	logger.l.Debug("", args...)
}

func (logger *Logger) Info(args ...any) {
	logger.l.Info("", args...)
}

func (logger *Logger) Warn(args ...any) {
	logger.l.Warn("", args...)
}

func (logger *Logger) Error(args ...any) {
	logger.l.Error("", args...)
}

func (logger *Logger) Fatal(args ...any) {
	logger.l.Error("", args...)
}
