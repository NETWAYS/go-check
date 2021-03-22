package check

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func ExampleConfig() {
	config := NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"

	_ = config.FlagSet.StringP("hostname", "H", "localhost", "Hostname to check")

	os.Args = []string{"check_example", "--help"}

	config.ParseArguments()

	log.Info("test")
	// Output: Usage of check_test
	//
	// Test Plugin
	//
	// Arguments:
	//   -H, --hostname string   Hostname to check (default "localhost")
	//   -t, --timeout int       Abort the check after n seconds (default 30)
	//   -d, --debug             Enable debug mode
	//   -v, --verbose           Enable verbose mode
	//   -V, --version           Print version and exit
	// would exit with code 3
}
