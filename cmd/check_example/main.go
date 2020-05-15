package main

import (
	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	defer check.CatchPanic()
	flags := check.NewFlags()
	flags.Name = "check_test"
	flags.Readme = `Test Plugin`
	flags.Version = "1.0.0"
	flags.Timeout = 10

	value := flags.Set.IntP("value", "", 10, "test value")
	warning := flags.Set.IntP("warning", "w", 20, "warning threshold")
	critical := flags.Set.IntP("critical", "c", 50, "critical threshold")

	flags.ParseArguments()

	log.Info("Start logging")

	time.Sleep(20 * time.Second)

	if *value > *critical {
		check.Exit(check.Critical, "value is %d", *value)
	} else if *value > *warning {
		check.Exit(check.Warning, "value is %d", *value)
	} else {
		check.Exit(check.OK, "value is %d", *value)
	}
}
