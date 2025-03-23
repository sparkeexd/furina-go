package utils

import (
	"os"

	"github.com/sanity-io/litter"
	"gorm.io/gorm/logger"
)

// Application environment.
var env = os.Getenv("ENV")

// Dump data structures to aid in debugging and testing.
// Prevents dumping in production environment if this function is left in code.
func Dump(message ...any) {
	if env == "development" {
		litter.Dump(message...)
	}
}

// Sets log level for the GORM logger.
// Development: Info
// Production: Error
func LogLevel() logger.LogLevel {
	logLevel := logger.Error
	if env == "development" {
		logLevel = logger.Info
	}

	return logLevel
}
