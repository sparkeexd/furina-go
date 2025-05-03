package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Custom slog logger.
type Logger struct {
	slog.Logger
}

// Creates a new slog logger depending on the environment.
// Development: Custom colorized logger.
// Production: JSON logger.
func NewLogger() *Logger {
	var logger *Logger

	env := os.Getenv("ENV")
	if env == "development" {
		opts := LoggerHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: LevelTrace,
			},
		}

		handler := NewLoggerHandler(os.Stderr, opts)
		logger = &Logger{*slog.New(handler)}
	} else {
		opts := slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				if attr.Key == slog.LevelKey {
					level := attr.Value.Any().(slog.Level)
					levelLabel, exists := Levels[level]
					if !exists {
						panic(fmt.Sprintf("Unknown slog level: %v", level))
					}

					attr.Value = slog.StringValue(levelLabel.name)
				}

				return attr
			},
		}

		logger = &Logger{*slog.New(slog.NewJSONHandler(os.Stderr, &opts))}
	}

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
