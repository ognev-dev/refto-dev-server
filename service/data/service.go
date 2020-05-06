package data

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
)

type Type string

const (
	GenericType Type = "generic"
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

func JSONBytesFromYAMLFile(fPath string) (data []byte, err error) {
	yamlData, err := ioutil.ReadFile(fPath)
	if err != nil {
		return
	}

	data, err = yaml.YAMLToJSON(yamlData)
	return
}

func TypeFromFilename(fPath string) (t Type) {
	_, fName := path.Split(fPath)
	nameParts := strings.Split(fName, ".")
	if len(nameParts) > 2 {
		t = Type(nameParts[len(nameParts)-2])
	}

	if t == "" {
		t = GenericType
	}

	return
}

func TypeFromSchemaFilename(fPath string) (t Type) {
	_, fName := path.Split(fPath)
	nameParts := strings.Split(fName, ".")
	if len(nameParts) > 2 {
		t = Type(nameParts[0])
	}

	if t == "" {
		t = GenericType
	}

	return
}
