package utils

import (
	"os"

	"github.com/sanity-io/litter"
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
