package main

import (
	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	defer check.CatchPanic()
	flags := check.NewFlags()
	flags.Name = "check_test"
	flags.Readme = `Test Plugin`
	flags.Version = "1.0.0"

	value := flags.Set.IntP("value", "t", 10, "test value")
	warning := flags.Set.IntP("warning", "w", 20, "warning threshold")
	critical := flags.Set.IntP("critical", "c", 50, "critical threshold")

	// value should be calculated

	flags.Parse(os.Args[1:])
	flags.SetupLogging()

	log.Info("Start logging")

	if *value > *critical {
		check.Exit(check.Critical, "value is %d", *value)
	} else if *value > *warning {
		check.Exit(check.Warning, "value is %d", *value)
	} else {
		check.Exit(check.OK, "value is %d", *value)
	}
}
