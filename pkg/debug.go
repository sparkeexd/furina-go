package pkg

import (
	"encoding/json"
	"log"
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

// Print JSON response returned from API.
func PrintJSON(response any, err error) {
	jsonString, _ := json.Marshal(response)

	data := make(map[string]interface{})
	json.Unmarshal(jsonString, &data)

	if err != nil {
		log.Println(err)
	} else {
		data, _ := json.MarshalIndent(data, "", "    ")
		log.Println(string(data))
	}
}
