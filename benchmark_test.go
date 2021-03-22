package check

import (
	"time"
)

func ExampleNewBenchmark() {
	bench := NewBenchmark("Start plugin")
	defer bench.DumpWhen(true /* flags.Debug */)

	time.Sleep(1 * time.Second)
	bench.Record("Connecting to service")
	time.Sleep(2 * time.Second)
	bench.Record("Connection open")
}

//noinspection GoBoolExpressions
func ExampleInitBenchmark() {
	debug := true /* flags.Debug */
	if debug {
		InitBenchmark("Start plugin")

		defer DumpBenchmarkWhen(debug /* flags.Debug */)
	}

	time.Sleep(1 * time.Second)
	RecordBenchmark("Connecting to service")
	time.Sleep(2 * time.Second)
	RecordBenchmark("Connection open")
}
