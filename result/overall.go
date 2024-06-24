// result tries to
package result

import (
	"errors"
	"fmt"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
)

// Overall is a singleton for a monitoring pluging that has several partial results (or sub-results)
//
// Design decisions: A check plugin has a single Overall (singleton),
// each partial thing which is tested, gets its own subcheck.
//
// The results of these may be relevant to the overall status in the end
// or not, e.g. if a plugin tries two different methods for something and
// one suffices, but one fails, the whole check might be OK and only the subcheck
// Warning or Critical.
type Overall struct {
	oks                int
	warnings           int
	criticals          int
	unknowns           int
	Summary            string
	stateSetExplicitly bool
	Outputs            []string // Deprecate this in a future version
	PartialResults     []PartialResult
}

// PartialResult represents a sub-result for an Overall struct
type PartialResult struct {
	Perfdata           perfdata.PerfdataList
	PartialResults     []PartialResult
	Output             string
	state              int  // Result state, either set explicitly or derived from partialResults
	defaultState       int  // Default result state, if no partial results are available and no state is set explicitly
	stateSetExplicitly bool // nolint: unused
	defaultStateSet    bool // nolint: unused
}

// Initializer for a PartialResult with "sane" defaults
// Notable default compared to the nil object: the default state is set to Unknown
func NewPartialResult() PartialResult {
	return PartialResult{
		stateSetExplicitly: false,
		defaultState:       check.Unknown,
	}
}

// String returns the status and output of the PartialResult
func (s *PartialResult) String() string {
	return fmt.Sprintf("[%s] %s", check.StatusText(s.GetStatus()), s.Output)
}

// Add adds a return state explicitly
//
// Hint: This will set stateSetExplicitly to true
func (o *Overall) Add(state int, output string) {
	switch state {
	case check.OK:
		o.oks++
	case check.Warning:
		o.warnings++
	case check.Critical:
		o.criticals++
	default:
		o.unknowns++
	}

	// TODO: Might be a bit obscure that the Add method also sets stateSetExplicitly
	o.stateSetExplicitly = true

	o.Outputs = append(o.Outputs, fmt.Sprintf("[%s] %s", check.StatusText(state), output))
}

// AddSubcheck adds a PartialResult to the Overall
func (o *Overall) AddSubcheck(subcheck PartialResult) {
	o.PartialResults = append(o.PartialResults, subcheck)
}

// AddSubcheck adds a PartialResult to the PartialResult
func (s *PartialResult) AddSubcheck(subcheck PartialResult) {
	s.PartialResults = append(s.PartialResults, subcheck)
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the Overall
func (o *Overall) GetStatus() int {
	if o.stateSetExplicitly {
		// nolint: gocritic
		if o.criticals > 0 {
			return check.Critical
		} else if o.unknowns > 0 {
			return check.Unknown
		} else if o.warnings > 0 {
			return check.Warning
		} else if o.oks > 0 {
			return check.OK
		} else {
			return check.Unknown
		}
	} else {
		// state not set explicitly!
		if len(o.PartialResults) == 0 {
			return check.Unknown
		}

		var (
			criticals int
			warnings  int
			oks       int
			unknowns  int
		)

		for _, sc := range o.PartialResults {
			switch sc.GetStatus() {
			case check.Critical:
				criticals++
			case check.Warning:
				warnings++
			case check.Unknown:
				unknowns++
			case check.OK:
				oks++
			}
		}

		if criticals > 0 {
			return check.Critical
		}

		if unknowns > 0 {
			return check.Unknown
		}

		if warnings > 0 {
			return check.Warning
		}

		if oks > 0 {
			return check.OK
		}

		return check.Unknown
	}
}

// GetSummary returns a text representation of the current state of the Overall
// nolint: funlen
func (o *Overall) GetSummary() string {
	if o.Summary != "" {
		return o.Summary
	}

	// Was the state set explicitly?
	if o.stateSetExplicitly {
		// Yes, so lets generate it from the sum of the overall states
		if o.criticals > 0 {
			o.Summary += fmt.Sprintf("critical=%d ", o.criticals)
		}

		if o.unknowns > 0 {
			o.Summary += fmt.Sprintf("unknown=%d ", o.unknowns)
		}

		if o.warnings > 0 {
			o.Summary += fmt.Sprintf("warning=%d ", o.warnings)
		}

		if o.oks > 0 {
			o.Summary += fmt.Sprintf("ok=%d ", o.oks)
		}

		if o.Summary == "" {
			o.Summary = "No status information"
			return o.Summary
		}
	}

	if !o.stateSetExplicitly {
		// No, so lets combine the partial ones
		if len(o.PartialResults) == 0 {
			// Oh, we actually don't have those either
			o.Summary = "No status information"
			return o.Summary
		}

		var (
			criticals int
			warnings  int
			oks       int
			unknowns  int
		)

		for _, sc := range o.PartialResults {
			switch sc.GetStatus() {
			case check.Critical:
				criticals++
			case check.Warning:
				warnings++
			case check.Unknown:
				unknowns++
			case check.OK:
				oks++
			}
		}

		if criticals > 0 {
			o.Summary += fmt.Sprintf("critical=%d ", criticals)
		}

		if unknowns > 0 {
			o.Summary += fmt.Sprintf("unknowns=%d ", unknowns)
		}

		if warnings > 0 {
			o.Summary += fmt.Sprintf("warning=%d ", warnings)
		}

		if oks > 0 {
			o.Summary += fmt.Sprintf("ok=%d ", oks)
		}
	}

	o.Summary = "states: " + strings.TrimSpace(o.Summary)

	return o.Summary
}

// GetOutput returns a text representation of the current outputs of the Overall
func (o *Overall) GetOutput() string {
	var output strings.Builder

	output.WriteString(o.GetSummary() + "\n")

	for _, extra := range o.Outputs {
		output.WriteString(extra + "\n")
	}

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

// getPerfdata returns all subsequent perfdata as a concatenated string
func (s *PartialResult) getPerfdata() string {
	var output strings.Builder

	if len(s.Perfdata.List) > 0 {
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

	prefix := strings.Repeat("  ", indentLevel)
	output.WriteString(prefix + "\\_ " + s.String() + "\n")

	if s.PartialResults != nil {
		for _, ss := range s.PartialResults {
			output.WriteString(ss.getOutput(indentLevel + 2))
		}
	}

	return output.String()
}

// SetDefaultState sets a new default state for a PartialResult
func (s *PartialResult) SetDefaultState(state int) error {
	if state < check.OK || state > check.Unknown {
		return errors.New("Default State is not a valid result state. Got " + fmt.Sprint(state) + " which is not valid")
	}

	s.defaultState = state
	s.defaultStateSet = true

	return nil
}

// SetState sets a state for a PartialResult
func (s *PartialResult) SetState(state int) error {
	if state < check.OK || state > check.Unknown {
		return errors.New("Default State is not a valid result state. Got " + fmt.Sprint(state) + " which is not valid")
	}

	s.state = state
	s.stateSetExplicitly = true

	return nil
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the PartialResult
// nolint: unused
func (s *PartialResult) GetStatus() int {
	if s.stateSetExplicitly {
		return s.state
	}

	if len(s.PartialResults) == 0 {
		if s.defaultStateSet {
			return s.defaultState
		}

		return check.Unknown
	}

	states := make([]int, len(s.PartialResults))

	for i := range s.PartialResults {
		states[i] = s.PartialResults[i].GetStatus()
	}

	return WorstState(states...)
}
