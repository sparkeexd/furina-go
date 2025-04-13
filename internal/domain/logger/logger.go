package logger

import (
	"context"
	"log/slog"
	"os"
)

// Custom slog logger.
type Logger struct {
	slog.Logger
}

// Creates a new slog logger with a custom handler that formats slog messages.
func NewLogger() *Logger {
	opts := LoggerHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: LevelTrace,
		},
	}

	handler := NewLoggerHandler(os.Stdout, opts)
	logger := &Logger{*slog.New(handler)}

	return logger
}

// Debug logs at [LevelTrace].
func (logger *Logger) Trace(msg string, args ...any) {
	logger.Log(context.Background(), LevelTrace, msg, args...)
}

// Debug logs at [LevelFatal].
func (logger *Logger) Fatal(msg string, args ...any) {
	logger.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}
