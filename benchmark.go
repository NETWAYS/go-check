package check

import (
	"fmt"
	"github.com/NETWAYS/go-check/convert"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// Benchmark records multiple events and provides functionality to dump the recorded events
type Benchmark struct {
	Events []*BenchmarkEvent
}

// BenchmarkEvent represents a single event during a benchmark with the time it occurred
type BenchmarkEvent struct {
	Time       *time.Time
	Offset     *time.Duration
	Message    string
	TotalAlloc uint64
	HeapAlloc  uint64
}

var ActiveBenchmark *Benchmark

// NewBenchmark starts a new benchmark and records the message as event
func NewBenchmark(message string) *Benchmark {
	b := &Benchmark{}
	b.Record(message)
	return b
}

// Initialize a Benchmark for the running program as static instance
func InitBenchmark(message string) {
	ActiveBenchmark = NewBenchmark(message)
}

// Recording timing and memory for the current program
func RecordBenchmark(message string) {
	if ActiveBenchmark != nil {
		ActiveBenchmark.Record(message)
	}
}

// Dump the Benchmark to stdout
func DumpBenchmark() {
	ActiveBenchmark.Dump()
}

// Dump the Benchmark when criteria is met, to be used with defer and a boolean variable
func DumpBenchmarkWhen(criteria bool) {
	ActiveBenchmark.DumpWhen(criteria)
}

// Record the current step of execution
//
// Use a short description, so that one knows at what stage the program is at. Time and offset will be calculated
// automatically.
func (b *Benchmark) Record(message string) {
	t := time.Now()

	// calculate duration since last event
	count := len(b.Events)
	var dur time.Duration
	if count > 0 {
		dur = t.Sub(*b.Events[count-1].Time)
	}

	// read memory info
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	// log the data for other use
	log.WithFields(log.Fields{
		"offset":      dur,
		"total_alloc": mem.TotalAlloc,
		"heap_alloc":  mem.HeapAlloc,
	}).Debug("Benchmark: ", message)

	e := &BenchmarkEvent{&t, &dur, message, mem.TotalAlloc, mem.HeapAlloc}
	b.Events = append(b.Events, e)
}

// Dump the benchmark to stdout
func (b *Benchmark) Dump() {
	fmt.Println("\n## Benchmark")

	// TODO: better table alignment
	fmt.Println("clock    | offset | total | heap | message")
	fmt.Println("---------|--------|-------|------|--------")

	for _, event := range b.Events {
		total := &convert.Bytesize{Data: event.TotalAlloc, Unit: "B"}
		heap := &convert.Bytesize{Data: event.HeapAlloc, Unit: "B"}

		fmt.Printf("%s | %.03f | %0.2f MB | %0.2f MB | %s\n",
			event.Time.Format("15:04:05"),
			event.Offset.Seconds(),
			total.ToMegabyte(),
			heap.ToMegabyte(),
			event.Message)
	}
}

// Dump the benchmark when a criteria is met
//
// Makes it simple to defer dumping by adding a boolean reference from a variable.
func (b *Benchmark) DumpWhen(criteria bool) {
	if criteria {
		b.Dump()
	}
}
