package main

import (
	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer check.CatchPanic()
	config := check.NewConfig()
	config.Name = "check_test"
	config.Readme = `Test Plugin`
	config.Version = "1.0.0"
	config.Timeout = 10

	value := config.FlagSet.IntP("value", "", 10, "test value")
	warning := config.FlagSet.IntP("warning", "w", 20, "warning threshold")
	critical := config.FlagSet.IntP("critical", "c", 50, "critical threshold")

	config.ParseArguments()

	log.Info("Start logging")

	// time.Sleep(20 * time.Second)

	if *value > *critical {
		check.Exitf(check.Critical, "value is %d", *value)
	} else if *value > *warning {
		check.Exitf(check.Warning, "value is %d", *value)
	} else {
		check.Exitf(check.OK, "value is %d", *value)
	}
}
