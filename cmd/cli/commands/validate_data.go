package commands

import (
	jsonschema "github.com/refto/server/service/json_schema"
)

func init() {
	add("validate", command{
		handler: validateData,
		help:    "Validates data against existing JSON schemas",
	})
}

func validateData(args ...string) (err error) {
	return jsonschema.Validate()
}
