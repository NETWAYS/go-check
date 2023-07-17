package check

import (
	"fmt"
	"os"
	"path"

	flag "github.com/spf13/pflag"
)

// Config represents a configuration for a monitoring plugin's CLI
type Config struct {
	// Name of the monitoring plugin
	Name string
	// README represents the help text for the CLI usage
	Readme string
	// Output for the --version flag
	Version string
	// Default for the --timeout flag
	Timeout int
	// Default for the --verbose flag
	Verbose bool
	// Default for the --debug flag
	Debug bool
	// Enable predefined --version output
	PrintVersion bool
	// Enable predefined default flags for the monitoring plugin
	DefaultFlags bool
	// Enable predefined default functions (e.g. Timeout handler) for the monitoring plugin
	DefaultHelper bool
	// Additional CLI flags for the monitoring plugin
	FlagSet *flag.FlagSet
}

// NewConfig returns a Config struct with some defaults
func NewConfig() *Config {
	c := &Config{}
	c.Name = path.Base(os.Args[0])

	c.FlagSet = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	c.FlagSet.SortFlags = false
	c.FlagSet.SetOutput(os.Stdout)

	c.FlagSet.Usage = func() {
		fmt.Printf("Usage of %s\n", c.Name)

		if c.Readme != "" {
			fmt.Println()
			fmt.Println(c.Readme)
		}

		fmt.Println()
		fmt.Println("Arguments:")
		c.FlagSet.PrintDefaults()
	}

	// set some defaults
	c.DefaultFlags = true
	c.Timeout = 30
	c.DefaultHelper = true

	return c
}

// ParseArguments parses the command line arguments given by os.Args
func (c *Config) ParseArguments() {
	c.ParseArray(os.Args[1:])
}

// ParseArray parses a list of command line arguments
func (c *Config) ParseArray(arguments []string) {
	if c.DefaultFlags {
		c.addDefaultFlags()
	}

	err := c.FlagSet.Parse(arguments)
	if err != nil {
		ExitError(err)
	}

	if c.PrintVersion {
		fmt.Println(c.Name, "version", c.Version)
		BaseExit(3)
	}

	if c.DefaultHelper {
		c.EnableTimeoutHandler()
	}
}

// addDefaultFlags adds various default flags to the monitoring plugin
func (c *Config) addDefaultFlags() {
	c.FlagSet.IntVarP(&c.Timeout, "timeout", "t", c.Timeout, "Abort the check after n seconds")
	c.FlagSet.BoolVarP(&c.Debug, "debug", "d", false, "Enable debug mode")
	c.FlagSet.BoolVarP(&c.Verbose, "verbose", "v", false, "Enable verbose mode")
	c.FlagSet.BoolVarP(&c.PrintVersion, "version", "V", false, "Print version and exit")

	c.DefaultFlags = false
}
