package check

func ExampleConfig() {
	config := NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"

	_ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

	config.ParseArguments()

	// Some checking should be done here

	Exitf(OK, "Everything is fine - answer=%d", 42)

	// Output: [OK] - Everything is fine - answer=42
	// would exit with code 0
}
