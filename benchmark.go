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
func NewBenchmark(message string) (b *Benchmark) {
	b = &Benchmark{}
	b.Record(message)

	return
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

	var dur time.Duration

	// calculate duration since last event
	count := len(b.Events)
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
		total := convert.BytesIEC(event.TotalAlloc)
		heap := convert.BytesIEC(event.HeapAlloc)

		fmt.Printf("%s | %.03f | %s | %s | %s\n",
			event.Time.Format("15:04:05"),
			event.Offset.Seconds(),
			total.HumanReadable(),
			heap.HumanReadable(),
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
