package jsonschema

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ognev-dev/bits/config"
	"github.com/ognev-dev/bits/service/data"
	"github.com/xeipuuv/gojsonschema"
)

type Schema struct {
	// I need Path only to display informative error messages
	Path   string
	Schema *gojsonschema.Schema
}

type File struct {
	// I need Path only to display informative error messages
	Path string
	Data []byte
}

type schemas map[data.Type]Schema
type bits map[data.Type][]File

type errFile struct {
	Path    string
	Message string
}

type errStack []string

func (e *errStack) Add(err error, wrapOpt ...string) {
	if len(wrapOpt) > 0 {
		err = errors.New(wrapOpt[0] + ": " + err.Error())
	}

	*e = append(*e, err.Error())
}
func (e errStack) Error() string {
	return strings.Join(e, "\n")
}

// Validate validates each YAML doc against JSON schema of it's type
func Validate() (err error) {
	schemasRepo := schemas{}
	bitsRepo := bits{}
	errs := errStack{}

	// collect schemas along with data (to not walk dirs seconds time)
	err = filepath.Walk(config.Get().Dir.Data, func(path string, f os.FileInfo, wErr error) (err error) {
		if wErr != nil {
			return wErr
		}

		if data.IsSchemaFile(f.Name()) {
			err = registerSchema(path, schemasRepo)
			if err != nil {
				errs.Add(err, "register schema: "+relPath(path))
			}
			return nil
		}

		// skip dirs, samples and non-yaml files
		if !data.IsDataFile(f) {
			return
		}

		jsonBytes, err := data.JSONBytesFromYAMLFile(path)
		if err != nil {
			errs.Add(err, relPath(path))
			return nil
		}

		t := data.TypeFromFilename(path)
		bitsByType, ok := bitsRepo[t]
		if !ok {
			bitsByType = []File{}
		}
		bitsByType = append(bitsByType, File{Path: path, Data: jsonBytes})
		bitsRepo[t] = bitsByType

		return nil
	})

	if len(errs) > 0 {
		err = errs
		return
	}

	// iterate over collected data and validate each
	for t, bitsByType := range bitsRepo {
		for _, v := range bitsByType {
			schema, ok := schemasRepo[t]
			if !ok {
				errs.Add(fmt.Errorf("schema of type '%s' is not exists (source '%s')", t, relPath(v.Path)))
				break
			}

			loader := gojsonschema.NewBytesLoader(v.Data)
			result, err := schema.Schema.Validate(loader)
			if err != nil {
				errs.Add(err, "validate "+relPath(v.Path))
				continue
			}

			if !result.Valid() {
				errs.Add(errors.New("validation failed for " + relPath(v.Path)))
				for _, e := range result.Errors() {
					errs.Add(errors.New("\t - " + e.String()))
				}
			}
		}
	}

	if len(errs) > 0 {
		err = errs
	}

	return
}

// registerSchema loads YAML schema from fPath converts it to JSON and adds it to the repo
// returns error if schema of given type already been registered
func registerSchema(fPath string, repo schemas) (err error) {
	t := data.TypeFromSchemaFilename(fPath)

	// check if this schema type already registered
	existing, ok := repo[t]
	if ok {
		err = fmt.Errorf(
			"schema type '%s' from '%s' already registered in '%s'",
			t, relPath(fPath), relPath(existing.Path),
		)
		return
	}

	jsonBytes, err := data.JSONBytesFromYAMLFile(fPath)
	if err != nil {
		return
	}

	sl := gojsonschema.NewSchemaLoader()
	loader := gojsonschema.NewBytesLoader(jsonBytes)
	schema, err := sl.Compile(loader)
	if err != nil {
		return
	}

	repo[t] = Schema{
		Path:   fPath,
		Schema: schema,
	}

	return
}

func relPath(fPath string) string {
	return strings.TrimPrefix(fPath, config.Get().Dir.Data)
}
