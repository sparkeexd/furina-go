package pkg

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sanity-io/litter"
)

// Dump data structures to aid in debugging and testing.
func Dump(message ...any) {
	env := os.Getenv("ENV")
	if env != "development" {
		return
	}

	litter.Dump(message...)
}

// Print JSON response returned from API.
func PrintJSON(response any, err error) {
	env := os.Getenv("ENV")
	if env != "development" {
		return
	}

	jsonString, _ := json.Marshal(response)

	data := make(map[string]any)
	json.Unmarshal(jsonString, &data)

	if err != nil {
		log.Println(err)
	} else {
		data, _ := json.MarshalIndent(data, "", "    ")
		log.Println(string(data))
	}
}
