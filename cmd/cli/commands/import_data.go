package commands

import dataimport "github.com/refto/server/service/data_import"

func init() {
	add("import", command{
		handler: importData,
		help:    "Imports data into database",
	})
}

func importData(args ...string) (err error) {
	return dataimport.Import()
}
