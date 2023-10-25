// result tries to
package result

import (
	"encoding/json"
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

type OverallOutput struct {
	MpiVersion     uint                  `json:"mpi_version"`
	Rc             int                   `json:"rc"`
	Output         string                `json:"output,omitempty"`
	PartialResults []PartialResultOutput `json:"partial_results,omitempty"`
	Perfdata       perfdata.PerfdataList `json:"perfata,omitempty"`
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

type PartialResultOutput struct {
	Rc             int                   `json:"rc"`
	Output         string                `json:"output,omitempty"`
	PartialResults []PartialResultOutput `json:"partial_results,omitempty"`
	Perfdata       perfdata.PerfdataList `json:"perfata,omitempty"`
}

// String returns the status and output of the PartialResult
func (pr *PartialResult) String() string {
	return fmt.Sprintf("[%s] %s", check.StatusText(pr.GetStatus()), pr.Output)
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
func (pr *PartialResult) AddSubcheck(subcheck PartialResult) {
	pr.PartialResults = append(pr.PartialResults, subcheck)
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
func (pr *PartialResult) getPerfdata() string {
	var output strings.Builder

	if len(pr.Perfdata) > 0 {
		output.WriteString(pr.Perfdata.String())
	}

	if pr.PartialResults != nil {
		for _, ss := range pr.PartialResults {
			output.WriteString(" " + ss.getPerfdata())
		}
	}

	return strings.TrimSpace(output.String())
}

// getOutput generates indented output for all subsequent PartialResults
func (pr *PartialResult) getOutput(indentLevel int) string {
	var output strings.Builder

	prefix := strings.Repeat("  ", indentLevel)
	output.WriteString(prefix + "\\_ " + pr.String() + "\n")

	if pr.PartialResults != nil {
		for _, ss := range pr.PartialResults {
			output.WriteString(ss.getOutput(indentLevel + 2))
		}
	}

	return output.String()
}

// SetDefaultState sets a new default state for a PartialResult
func (pr *PartialResult) SetDefaultState(state int) error {
	if state < check.OK || state > check.Unknown {
		return errors.New("Default State is not a valid result state. Got " + fmt.Sprint(state) + " which is not valid")
	}

	pr.defaultState = state
	pr.defaultStateSet = true

	return nil
}

// SetState sets a state for a PartialResult
func (pr *PartialResult) SetState(state int) error {
	if state < check.OK || state > check.Unknown {
		return errors.New("Default State is not a valid result state. Got " + fmt.Sprint(state) + " which is not valid")
	}

	pr.state = state
	pr.stateSetExplicitly = true

	return nil
}

// GetStatus returns the current state (ok, warning, critical, unknown) of the PartialResult
// nolint: unused
func (pr *PartialResult) GetStatus() int {
	if pr.stateSetExplicitly {
		return pr.state
	}

	if len(pr.PartialResults) == 0 {
		if pr.defaultStateSet {
			return pr.defaultState
		}

		return check.Unknown
	}

	states := make([]int, len(pr.PartialResults))

	for i := range pr.PartialResults {
		states[i] = pr.PartialResults[i].state
	}

	return WorstState(states...)
}

func (pr *PartialResult) convertToOutput() PartialResultOutput {
	result := PartialResultOutput{}
	result.Output = pr.Output
	result.Perfdata = pr.Perfdata
	result.Rc = pr.GetStatus()

	if len(pr.PartialResults) != 0 {
		for i := range pr.PartialResults {
			tmp := pr.PartialResults[i].convertToOutput()
			result.PartialResults = append(result.PartialResults, tmp)
		}
	}

	return result
}

func (o *Overall) convertToOutput(version uint) OverallOutput {
	result := OverallOutput{}
	result.Output = o.Summary
	result.Rc = o.GetStatus()
	result.MpiVersion = version

	if len(o.PartialResults) != 0 {
		for i := range o.PartialResults {
			tmp := o.PartialResults[i].convertToOutput()
			result.PartialResults = append(result.PartialResults, tmp)
		}
	}

	return result
}

func (o *Overall) GetMpiOutput(version uint) []byte {
	oo := o.convertToOutput(version)

	result, err := json.Marshal(oo)
	if err != nil {
		return []byte{}
	}

	return result
}
