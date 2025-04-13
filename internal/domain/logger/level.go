package logger

import "log/slog"

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorDarkGray = 90

	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var (
	// Maps slog levels to ANSI color codes.
	Levels = map[slog.Level]levelConfig{
		LevelTrace: {
			name:       "TRACE",
			levelColor: colorBlue,
			keyColor:   colorCyan,
			valueColor: colorDarkGray,
		},
		slog.LevelDebug: {
			name:       "DEBUG",
			levelColor: colorMagenta,
			keyColor:   colorCyan,
			valueColor: colorDarkGray,
		},
		slog.LevelInfo: {
			name:       "INFO",
			levelColor: colorGreen,
			keyColor:   colorCyan,
			valueColor: colorDarkGray,
		},
		slog.LevelWarn: {
			name:       "WARN",
			levelColor: colorYellow,
			keyColor:   colorCyan,
			valueColor: colorYellow,
		},
		slog.LevelError: {
			name:       "ERROR",
			levelColor: colorRed,
			keyColor:   colorCyan,
			valueColor: colorRed,
		},
		LevelFatal: {
			name:       "FATAL",
			levelColor: colorRed,
			keyColor:   colorCyan,
			valueColor: colorRed,
		},
	}
)

// Color configuration for a log level.
type levelConfig struct {
	name       string
	levelColor int
	keyColor   int
	valueColor int
}
