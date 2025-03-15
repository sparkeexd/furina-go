package utils

import (
	"os"

	"github.com/sanity-io/litter"
)

// Dump data structures to aid in debugging and testing.
// Prevents dumping in production environment if this function is left in code.
func Dump(message ...any) {
	env := os.Getenv("ENV")

	if env == "dev" {
		litter.Dump(message...)
	}
}
