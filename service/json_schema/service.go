package jsonschema

import (
	"github.com/xeipuuv/gojsonschema"
)

func Validate(dataPath string) (err error) {
	schemas := map[string]gojsonschema.JSONLoader{
		"generic": nil,
	}
	// validate
	schemaLoader := gojsonschema.NewStringLoader("file:///home/me/schema.json")
	documentLoader := gojsonschema.NewReferenceLoader("file:///home/me/document.json")

	return
}
