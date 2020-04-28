package commands

func init() {
	add("seed", "sd", command{
		handler: seed,
		help:    `Seeds database with test data. Should be run on empty database`,
	})
}

func seed(args ...string) (err error) {

	// todo

	return
}
