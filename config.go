package check

import (
	"fmt"
	"os"
	"path"

	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Name          string
	Readme        string
	Version       string
	Timeout       int
	Verbose       bool
	Debug         bool
	PrintVersion  bool
	DefaultFlags  bool
	DefaultHelper bool
	FlagSet       *flag.FlagSet
}

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

func (c *Config) ParseArguments() {
	c.ParseArray(os.Args[1:])
}

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
		c.SetupLogging()
		c.EnableTimeoutHandler()
	}
}

func (c *Config) addDefaultFlags() {
	c.FlagSet.IntVarP(&c.Timeout, "timeout", "t", c.Timeout, "Abort the check after n seconds")
	c.FlagSet.BoolVarP(&c.Debug, "debug", "d", false, "Enable debug mode")
	c.FlagSet.BoolVarP(&c.Verbose, "verbose", "v", false, "Enable verbose mode")
	c.FlagSet.BoolVarP(&c.PrintVersion, "version", "V", false, "Print version and exit")

	c.DefaultFlags = false
}

func (c *Config) SetupLogging() {
	if c.Debug {
		log.SetLevel(log.DebugLevel)
	} else if c.Verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
