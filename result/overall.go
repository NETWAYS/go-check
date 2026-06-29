// Package result provides types and functions to organize results in a check plugin
package result

import (
	"fmt"
	"strings"
	"sync"

	"github.com/NETWAYS/go-check"
)

// The "width" of the indentation which is added on every level
const indentationOffset = 2

// statusCount is used to count the overall states
type statusCount struct {
	OK       int
	Warning  int
	Critical int
	Unknown  int
}

// Overall is a singleton for a monitoring plugin that has several partial results (or sub-results)
//
// Design decisions: A check plugin has a single Overall (singleton),
// each partial thing which is tested, gets its own subcheck.
//
// The results of these may be relevant to the overall status in the end
// or not, e.g. if a plugin tries two different methods for something and
// one suffices, but one fails, the whole check might be OK and only the subcheck
// Warning or Critical.
type Overall struct {
	// default summary (first line of output) if everything is ok. Has to be set in a plugin
	OKSummary string
	// The results that are associated with this overall
	PartialResults []*PartialResult

	// We use a Mutex to make sure PartialResults can be added and evaluated concurrently
	mu sync.RWMutex
}

// Add adds a return state explicitly.
// Add is concurrency-safe
func (o *Overall) Add(state check.Status, output string) {
	var result PartialResult
	result.SetState(state)
	result.SetOutput(output)
	o.AddSubcheck(&result)
}

// AddSubcheck adds a PartialResult to the Overall.
// AddSubcheck is concurrency-safe
func (o *Overall) AddSubcheck(subcheck *PartialResult) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.PartialResults = append(o.PartialResults, subcheck)
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the Overall.
// GetStatus is concurrency-safe
func (o *Overall) GetStatus() check.Status {
	o.mu.RLock()
	defer o.mu.RUnlock()

	statuses := o.getStatusCount()

	if statuses.Critical > 0 {
		return check.Critical
	}

	if statuses.Unknown > 0 {
		return check.Unknown
	}

	if statuses.Warning > 0 {
		return check.Warning
	}

	if statuses.OK > 0 {
		return check.OK
	}

	return check.Unknown
}

// GetOutput returns a text representation of the current outputs of the Overall.
// GetOutput is concurrency-safe
func (o *Overall) GetOutput() string {
	o.mu.RLock()
	defer o.mu.RUnlock()

	var output strings.Builder

	output.WriteString(o.getSummary() + "\n")

	if o.PartialResults != nil {
		var pdata strings.Builder

		// Generate indeted output and perfdata for all partialResults
		for i := range o.PartialResults {
			output.WriteString(strings.ReplaceAll(o.PartialResults[i].getOutput(0), check.PerfdataSeparatorSymbol, " "))
			pdata.WriteString(" " + o.PartialResults[i].getPerfdata())
		}

		pdataString := strings.Trim(pdata.String(), " ")

		if len(pdataString) > 0 {
			output.WriteString(check.PerfdataSeparatorSymbol + pdataString + "\n")
		}
	}

	return output.String()
}

func (o *Overall) getStatusCount() statusCount {
	result := statusCount{
		OK:       0,
		Warning:  0,
		Critical: 0,
		Unknown:  0,
	}

	if len(o.PartialResults) == 0 {
		return result
	}

	for _, sc := range o.PartialResults {
		switch sc.GetStatus() {
		case check.Critical:
			result.Critical++
		case check.Warning:
			result.Warning++
		case check.Unknown:
			result.Unknown++
		case check.OK:
			result.OK++
		}
	}

	return result
}

// GetSummary returns a text representation of the current state of the Overall
func (o *Overall) getSummary() string {
	checkState := o.GetStatus()

	if checkState == check.OK && o.OKSummary != "" {
		return strings.ReplaceAll(o.OKSummary, check.PerfdataSeparatorSymbol, " ")
	}

	if len(o.PartialResults) == 0 {
		// Oh, we actually don't have those either
		return "No status information"
	}

	if checkState == check.OK {
		return o.getGenericSummary()
	}

	result := ""
	worstState := check.OK

	// Get the worst non-ok PartialResults output
	for _, partRes := range o.PartialResults {
		if check.Compare(worstState, partRes.GetStatus()) > 0 {
			result = partRes.getPartialResultFailedOutput()
			worstState = partRes.GetStatus()
		}
	}

	if result == "" {
		// No output in PartialResults thus we generate the generic summary
		result = o.getGenericSummary()
	}

	return result
}

func (o *Overall) getGenericSummary() string {
	stats := o.getStatusCount()
	result := ""

	if stats.Critical > 0 {
		result += fmt.Sprintf("critical=%d ", stats.Critical)
	}

	if stats.Unknown > 0 {
		result += fmt.Sprintf("unknown=%d ", stats.Unknown)
	}

	if stats.Warning > 0 {
		result += fmt.Sprintf("warning=%d ", stats.Warning)
	}

	if stats.OK > 0 {
		result += fmt.Sprintf("ok=%d ", stats.OK)
	}

	result = "states: " + strings.TrimSpace(result)

	return result
}
