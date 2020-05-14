package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
)

func initRand() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 100
	return rand.Intn(max-min+1) + min
}

func main() {
	flags := check.NewFlags()
	flags.Name = "check_test"
	flags.Readme = `Test Plugin`
	flags.Version = "1.0.0"

	warning := flags.Set.IntP("warning", "w", 20, "warning threshold")
	critical := flags.Set.IntP("critical", "c", 50, "critical threshold")

	ranNum := initRand()

	flags.Parse(os.Args[1:])
	flags.SetupLogging()

	log.Info("Start logging")

	if *critical < ranNum {
		check.Exit(check.Critical, "value is %d", ranNum)
	} else if *warning < ranNum {
		check.Exit(check.Warning, "value is %d", ranNum)
	} else {
		check.Exit(check.OK, "value is %d", ranNum)
	}
}
