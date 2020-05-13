package check

import (
	log "github.com/sirupsen/logrus"
)

func ExampleFlags() {
	flags := NewFlags()
	flags.Name = "check_test"
	flags.Readme = `Test Plugin`
	flags.Version = "1.0.0"

	_ = flags.Set.StringP("hostname", "H", "localhost", "Hostname to check")

	// flags.Parse(os.Args[1:])
	flags.Parse([]string{"--help"})
	flags.SetupLogging()

	log.Info("test")
	// Output: Usage of check_test
	//
	// Test Plugin
	//
	// Arguments:
	// would exit with code 3
}
