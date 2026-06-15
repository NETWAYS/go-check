package result

import (
	"fmt"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
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
	OKSummary      string
	PartialResults []PartialResult
}

// Add adds a return state explicitly
func (o *Overall) Add(state check.Status, output string) {
	var result PartialResult
	result.SetState(state)
	result.Output = output
	o.AddSubcheck(result)
}

// AddSubcheck adds a PartialResult to the Overall
func (o *Overall) AddSubcheck(subcheck PartialResult) {
	o.PartialResults = append(o.PartialResults, subcheck)
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the Overall
func (o *Overall) GetStatus() check.Status {
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

// GetOutput returns a text representation of the current outputs of the Overall
func (o *Overall) GetOutput() string {
	var output strings.Builder

	output.WriteString(o.getSummary() + "\n")

	if o.PartialResults != nil {
		var pdata strings.Builder

		// Generate indeted output and perfdata for all partialResults
		for i := range o.PartialResults {
			output.WriteString(o.PartialResults[i].getOutput(0))
			pdata.WriteString(" " + o.PartialResults[i].getPerfdata())
		}

		pdataString := strings.Trim(pdata.String(), " ")

		if len(pdataString) > 0 {
			output.WriteString("|" + pdataString + "\n")
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

// PartialResult represents a sub-result for an Overall struct
type PartialResult struct {
	Perfdata       perfdata.PerfdataList
	PartialResults []PartialResult
	Output         string

	// Result state, either set explicitly or derived from partialResults
	state check.Status
	// Default result state, if no partial results are available and no state is set explicitly
	defaultState check.Status

	// stateSetExplicitly indicates that SetState was called directly. When true,
	// GetStatus returns s.state unconditionally, bypassing PartialResults entirely.
	stateSetExplicitly bool
	// defaultStateSetExplicitly indicates that SetDefaultState was called. When true
	// and no PartialResults exist and no explicit state is set, GetStatus returns
	// s.defaultState instead of check.Unknown.
	defaultStateSetExplicitly bool
}

// NewPartialResult initializer with defaults. It is recommended to use NewPartialResult.
// The default compared to the nil object is the default state is set to Unknown.
func NewPartialResult() PartialResult {
	return PartialResult{
		stateSetExplicitly: false,
		defaultState:       check.Unknown,
	}
}

// AddSubcheck adds a PartialResult to the PartialResult
func (s *PartialResult) AddSubcheck(subcheck PartialResult) {
	s.PartialResults = append(s.PartialResults, subcheck)
}

// String returns the status and output of the PartialResult
func (s *PartialResult) String() string {
	return fmt.Sprintf("[%s] %s", s.GetStatus(), s.Output)
}

// SetDefaultState sets a new default state for a PartialResult
func (s *PartialResult) SetDefaultState(state check.Status) {
	s.defaultState = state
	s.defaultStateSetExplicitly = true
}

// SetState sets a state for a PartialResult
func (s *PartialResult) SetState(state check.Status) {
	s.state = state
	s.stateSetExplicitly = true
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the PartialResult
func (s *PartialResult) GetStatus() check.Status {
	if s.stateSetExplicitly {
		return s.state
	}

	if len(s.PartialResults) == 0 {
		if s.defaultStateSetExplicitly {
			return s.defaultState
		}

		return check.Unknown
	}

	states := make([]check.Status, len(s.PartialResults))

	for i := range s.PartialResults {
		states[i] = s.PartialResults[i].GetStatus()
	}

	return WorstState(states...)
}

// getPerfdata returns all subsequent perfdata as a concatenated string
func (s *PartialResult) getPerfdata() string {
	var output strings.Builder

	if len(s.Perfdata) > 0 {
		output.WriteString(s.Perfdata.String())
	}

	if s.PartialResults != nil {
		for _, ss := range s.PartialResults {
			output.WriteString(" " + ss.getPerfdata())
		}
	}

	return strings.TrimSpace(output.String())
}

// getOutput generates indented output for all subsequent PartialResults
func (s *PartialResult) getOutput(indentLevel int) string {
	var output strings.Builder
	// The final result will look like this:
	// [OK] Overall is OK
	// \_ [OK] My PartialResult
	output.WriteString(strings.Repeat("  ", indentLevel) + "\\_ " + s.String() + "\n")

	if s.PartialResults != nil {
		for _, ss := range s.PartialResults {
			output.WriteString(ss.getOutput(indentLevel + indentationOffset))
		}
	}

	return output.String()
}

// GetSummary returns a text representation of the current state of the Overall
func (o *Overall) getSummary() string {
	checkState := o.GetStatus()

	if checkState == check.OK && o.OKSummary != "" {
		return o.OKSummary
	}

	if len(o.PartialResults) == 0 {
		// Oh, we actually don't have those either
		return "No status information"
	}

	if checkState == check.OK {
		return o.getGenericSummary()
	}

	// Non-OK result
	// get worst-first non-ok PartialResults output
	result := ""
	worstState := check.OK

	for _, partRes := range o.PartialResults {
		if check.Compare(worstState, partRes.GetStatus()) > 0 {
			result = partRes.getPartialResultFailedOutput()
		}
	}

	if result == "" {
		// No output in PartialResults ...
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

func (s *PartialResult) getPartialResultFailedOutput() string {
	if len(s.PartialResults) == 0 {
		// this is a leave node
		return s.Output
	}

	worstState := check.OK
	result := ""

	for _, partRes := range s.PartialResults {
		if check.Compare(worstState, partRes.GetStatus()) > 0 {
			result = partRes.getPartialResultFailedOutput()
		}
	}

	return result
}
