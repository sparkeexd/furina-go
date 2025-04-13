package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
)

// Options for the custom logger handler.
type LoggerHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

// Custom handler that formats slog messages with colourized output.
type LoggerHandler struct {
	slog.Handler
	logger *log.Logger
}

// Create a new logger handler for slog.
func NewLoggerHandler(out io.Writer, opts LoggerHandlerOptions) *LoggerHandler {
	handler := &LoggerHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		logger:  log.New(out, "", 0),
	}

	return handler
}

// Handle the slog record with colourized output.
func (handler *LoggerHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := make(map[string]any, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		fields[attr.Key] = attr.Value.Any()
		return true
	})

	timeStr := handler.colorize(record.Time.Format("2006-01-02 15:04:05"), colorDarkGray)
	level := handler.colorize(Levels[record.Level].name, Levels[record.Level].levelColor)
	message := handler.colorize(record.Message, colorWhite)
	fieldsStr := handler.formatFields(record.Level, fields)

	handler.logger.Printf("%s | %s | %s %s", timeStr, level, message, fieldsStr)
	return nil
}

// Formats the message with the specified ANSI color code.
func (handler *LoggerHandler) colorize(value any, color int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", color, value)
}

// Formats the fields map into a "key=value" string.
func (handler *LoggerHandler) formatFields(level slog.Level, fields map[string]any) string {
	var result string
	for key, value := range fields {
		result += fmt.Sprintf("%s=%v ", handler.colorize(key, Levels[level].keyColor), handler.colorize(value, Levels[level].valueColor))
	}
	return strings.TrimSpace(result)
}
