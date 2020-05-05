package data

import (
	"os"
	"strings"
)

const (
	YAMLExt   = ".yaml"
	sampleExt = ".sample.yaml"
	schemaExt = ".schema.yaml"
)

func IsSampleFile(fName string) bool {
	return strings.HasSuffix(fName, sampleExt)
}

func IsSchemaFile(fName string) bool {
	return strings.HasSuffix(fName, schemaExt)
}

func IsYamlFile(fName string) bool {
	return strings.HasSuffix(fName, YAMLExt)
}

func IsDataFile(f os.FileInfo) (ok bool) {
	if f.IsDir() {
		return false
	}
	if IsSchemaFile(f.Name()) {
		return false
	}
	if IsSampleFile(f.Name()) {
		return false
	}
	if IsYamlFile(f.Name()) {
		return true
	}

	return false
}
