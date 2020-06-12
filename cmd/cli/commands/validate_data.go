package commands

import (
	"github.com/refto/server/config"
	jsonschema "github.com/refto/server/service/json_schema"
)

func init() {
	add("validate", command{
		handler: validateData,
		help:    "Validates data against existing JSON schemas",
	})
}

func validateData(args ...string) (err error) {
	dirPath := config.Get().Dir.Data
	if len(args) > 0 {
		dirPath = args[0]
	}

	return jsonschema.Validate(dirPath)
}
