package check

import (
	"fmt"
	"os"
	"path"

	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
)

type Flags struct {
	Name         string
	Readme       string
	Version      string
	Verbose      bool
	Debug        bool
	PrintVersion bool
	Set          *flag.FlagSet
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

	flagSet.BoolVarP(&flags.Debug, "debug", "d", false, "Enable debug mode")
	flagSet.BoolVarP(&flags.Verbose, "verbose", "v", false, "Enable verbose mode")
	flagSet.BoolVarP(&flags.PrintVersion, "version", "V", false, "Print version and exit")

	flags.Set = flagSet

	return flags
}

func (f *Flags) Parse(arguments []string) {
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
