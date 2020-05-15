package check

import (
	"fmt"
	"os"
	"path"

	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

type Flags struct {
	Name          string
	Readme        string
	Version       string
	Timeout       int
	Verbose       bool
	Debug         bool
	PrintVersion  bool
	DefaultFlags  bool
	DefaultHelper bool
	Set           *flag.FlagSet
}

func NewFlags() *Flags {
	flags := &Flags{}
	flags.Name = path.Base(os.Args[0])

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SortFlags = false

	flagSet.Usage = func() {
		fmt.Printf("Usage of %s\n", flags.Name)
		if flags.Readme != "" {
			fmt.Println()
			fmt.Println(flags.Readme)
		}
		fmt.Println()
		fmt.Println("Arguments:")
		flagSet.PrintDefaults()
	}

	flags.Set = flagSet

	// set some defaults
	flags.DefaultFlags = true
	flags.Timeout = 30
	flags.DefaultHelper = true

	return flags
}

func (f *Flags) ParseArguments() {
	f.ParseArray(os.Args[1:])
}

func (f *Flags) ParseArray(arguments []string) {
	if f.DefaultFlags {
		f.addDefaultFlags()
	}

	err := f.Set.Parse(arguments)
	if err != nil {
		if err != flag.ErrHelp {
			ExitError(err)
		}
		BaseExit(3)
	}

	if f.PrintVersion {
		fmt.Println(f.Name, "version", f.Version)
		BaseExit(3)
	}

	if f.DefaultHelper {
		f.SetupLogging()
		f.EnableTimeoutHandler()
	}
}

func (f *Flags) addDefaultFlags() {
	f.Set.IntVarP(&f.Timeout, "timeout", "t", f.Timeout, "Abort the check after n seconds")
	f.Set.BoolVarP(&f.Debug, "debug", "d", false, "Enable debug mode")
	f.Set.BoolVarP(&f.Verbose, "verbose", "v", false, "Enable verbose mode")
	f.Set.BoolVarP(&f.PrintVersion, "version", "V", false, "Print version and exit")

	f.DefaultFlags = false
}

func (f *Flags) SetupLogging() {
	if f.Debug {
		log.SetLevel(log.DebugLevel)
	} else if f.Verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}
