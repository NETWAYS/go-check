package result

import (
	"fmt"
	"strings"
	"sync"

	"github.com/NETWAYS/go-check"
)

// PartialResult represents a sub-result for an Overall struct.
// Note that, a PartialResult must not be used on its own but always with an Overall.
type PartialResult struct {
	perfdata       check.PerfdataList
	partialResults []*PartialResult
	output         string

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

	mu sync.RWMutex
}

// NewPartialResult initializer with defaults. It is recommended to use NewPartialResult.
// The default compared to the nil object is the default state is set to Unknown.
func NewPartialResult() *PartialResult {
	return &PartialResult{
		stateSetExplicitly: false,
		defaultState:       check.Unknown,
	}
}

// AddSubcheck adds a PartialResult to the PartialResult
func (s *PartialResult) AddSubcheck(subcheck *PartialResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.partialResults = append(s.partialResults, subcheck)
}

// AddPerfdata adds a Perfdata point to the PartialResult
func (s *PartialResult) AddPerfdata(perfdata *check.Perfdata) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.perfdata.Add(perfdata)
}

// String returns the status and output of the PartialResult
func (s *PartialResult) String() string {
	return fmt.Sprintf("[%s] %s", s.GetStatus(), strings.ReplaceAll(s.output, check.PerfdataSeparatorSymbol, " "))
}

// SetDefaultState sets a new default state for a PartialResult
func (s *PartialResult) SetDefaultState(state check.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.defaultState = state
	s.defaultStateSetExplicitly = true
}

// SetState sets a state for a PartialResult
func (s *PartialResult) SetState(state check.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = state
	s.stateSetExplicitly = true
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the PartialResult
func (s *PartialResult) GetStatus() check.Status {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.stateSetExplicitly {
		return s.state
	}

	if len(s.partialResults) == 0 {
		if s.defaultStateSetExplicitly {
			return s.defaultState
		}

		return check.Unknown
	}

	states := make([]check.Status, len(s.partialResults))

	for i := range s.partialResults {
		states[i] = s.partialResults[i].GetStatus()
	}

	return check.WorstState(states...)
}

// SetOutput sets the output of this PartialResult to the given string
func (s *PartialResult) SetOutput(output string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.output = output
}

// getPerfdata returns all subsequent perfdata as a concatenated string
func (s *PartialResult) getPerfdata() string {
	var output strings.Builder

	if len(s.perfdata) > 0 {
		output.WriteString(s.perfdata.String())
	}

	if s.partialResults != nil {
		for _, ss := range s.partialResults {
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

	if s.partialResults != nil {
		for _, ss := range s.partialResults {
			output.WriteString(ss.getOutput(indentLevel + indentationOffset))
		}
	}

	return output.String()
}

func (s *PartialResult) getPartialResultFailedOutput() string {
	if len(s.partialResults) == 0 {
		// this is a leave node
		return s.output
	}

	result := ""
	worstState := check.OK

	// Get the worst non-ok PartialResults output
	for _, partRes := range s.partialResults {
		if check.Compare(worstState, partRes.GetStatus()) > 0 {
			result = partRes.getPartialResultFailedOutput()
			worstState = partRes.GetStatus()
		}
	}

	return result
}
