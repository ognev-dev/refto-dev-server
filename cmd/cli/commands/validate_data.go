package commands

import (
	"fmt"
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
	dataPath := config.Get().Dir.Data
	if len(args) > 0 {
		dataPath = args[0]
	}

	fmt.Println("Checking data at " + dataPath)
	resp, err := jsonschema.Validate(dataPath)
	if err != nil {
		errs, ok := err.(jsonschema.ErrStack)
		if !ok {
			fmt.Println("ERROR:", err.Error())
			return
		}
		fmt.Println("ERRORS:")
		for _, e := range errs {
			fmt.Println(" -" + e)
		}
		return
	}

	var hasWarning bool
	if resp.SchemaCount == 0 {
		fmt.Println("[WARNING] No schemas found in " + dataPath + ", most likely this is not a data directory")
		hasWarning = true
	}
	if resp.DataCount == 0 {
		fmt.Println("[WARNING] No data found in " + dataPath + ", most likely this is not a data directory")
		hasWarning = true
	}

	if hasWarning {
		fmt.Println("Validation is not successful due to warnings")
		return
	}

	fmt.Println("SUCCESS!")
	fmt.Println("Schemas found:", resp.SchemaCount)
	fmt.Println("Data validated:", resp.DataCount)
	fmt.Println("Data by type:")
	for t, v := range resp.DataCountByType {
		fmt.Println(" -", t+":", v)
	}

	return
}
