package main

import (
	"fmt"
	"log"

	"github.com/NETWAYS/go-check"
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

	if config.Debug {
		log.Println("Start logging")
	}
	// time.Sleep(20 * time.Second)

	if *value > *critical {
		check.Exit(check.Critical, fmt.Sprintf("value is %d", *value))
	} else if *value > *warning {
		check.Exit(check.Warning, fmt.Sprintf("value is %d", *value))
	} else {
		check.Exit(check.OK, fmt.Sprintf("value is %d", *value))
	}
}
