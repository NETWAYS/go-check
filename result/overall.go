// result tries to
package result

import (
	"fmt"
	"strings"

	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
)

// So, this is the idea:
// A check plugin has a single Overall (singleton)
// Each partial thing which is tested, gets its own subcheck
// The results of these may be relevant to the overall status in the end
// or not, e.g. if a plugin tries two different methods for something and
// one suffices, but one fails, the whole check might be OK and only the subcheck
// Warning or Critical.
type Overall struct {
	oks                 int
	warnings            int
	criticals           int
	unknowns            int
	Summary             string
	stateSetExplicitely bool
	Outputs             []string // Deprecate this in a future version
	partialResults      []PartialResult
}

type PartialResult struct {
	State               int
	Output              string
	stateSetExplicitely bool // nolint: unused
	Perfdata            perfdata.PerfdataList
	partialResults      []PartialResult
}

func (s *PartialResult) String() string {
	return fmt.Sprintf("[%s] %s", check.StatusText(s.State), s.Output)
}

// Deprecated: Will be removed in a future version, use Add() instead
func (o *Overall) AddOK(output string) {
	o.Add(check.OK, output)
}

// Deprecated: Will be removed in a future version, use Add() instead
func (o *Overall) AddWarning(output string) {
	o.Add(check.Warning, output)
}

// Deprecated: Will be removed in a future version, use Add() instead
func (o *Overall) AddCritical(output string) {
	o.Add(check.Critical, output)
}

// Deprecated: Will be removed in a future version, use Add() instead
func (o *Overall) AddUnknown(output string) {
	o.Add(check.Unknown, output)
}

// Add State explicitely
// Hint: This will set stateSetExplicitely to true
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

	// TODO: Might be a bit obscure that the Add method also sets stateSetExplicitely
	o.stateSetExplicitely = true

	o.Outputs = append(o.Outputs, fmt.Sprintf("[%s] %s", check.StatusText(state), output))
}

func (o *Overall) AddSubcheck(subcheck PartialResult) {
	o.partialResults = append(o.partialResults, subcheck)
}

func (o *PartialResult) AddSubcheck(subcheck PartialResult) {
	o.partialResults = append(o.partialResults, subcheck)
}

func (o *Overall) GetStatus() int {
	if o.stateSetExplicitely {
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
		// state set explicitely!

		if len(o.partialResults) == 0 {
			return check.Unknown
		}

		var (
			criticals int
			warnings  int
			oks       int
			unknowns  int
		)

		for _, sc := range o.partialResults {
			switch sc.State {
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

// nolint: funlen
func (o *Overall) GetSummary() string {
	if o.Summary != "" {
		return o.Summary
	}

	// Was the state set explicitely?
	if o.stateSetExplicitely {
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

	if !o.stateSetExplicitely {
		// No, so lets combine the partial ones
		if len(o.partialResults) == 0 {
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

		for _, sc := range o.partialResults {
			switch sc.State {
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

func (o *Overall) GetOutput() string {
	var output strings.Builder

	output.WriteString(o.GetSummary() + "\n")

	for _, extra := range o.Outputs {
		output.WriteString(extra + "\n")
	}

	if o.partialResults != nil {
		var pdata strings.Builder

		// Generate indeted output and perfdata for all partialResults
		for i := range o.partialResults {
			output.WriteString(o.partialResults[i].getOutput(0))
			pdata.WriteString(" " + o.partialResults[i].getPerfdata())
		}

		pdata_string := strings.Trim(pdata.String(), " ")

		if len(pdata_string) > 0 {
			output.WriteString("|" + pdata_string + "\n")
		}
	}

	return output.String()
}

// Returns all subsequent perfdata as a concatenated string
func (s *PartialResult) getPerfdata() string {
	var output strings.Builder

	if len(s.Perfdata) > 0 {
		output.WriteString(s.Perfdata.String())
	}

	if s.partialResults != nil {
		for _, ss := range s.partialResults {
			output.WriteString(ss.getPerfdata())
		}
	}

	return output.String()
}

// Generates indented output for all subsequent PartialResults
func (s *PartialResult) getOutput(indent_level int) string {
	var output strings.Builder

	prefix := strings.Repeat("  ", indent_level)
	output.WriteString(prefix + "\\_ " + s.String() + "\n")

	if s.partialResults != nil {
		for _, ss := range s.partialResults {
			output.WriteString(ss.getOutput(indent_level + 2))
		}
	}

	return output.String()
}

// nolint: unused
func (s *PartialResult) getState() int {
	if s.stateSetExplicitely {
		return s.State
	}

	if len(s.partialResults) == 0 {
		return check.Unknown
	}

	states := make([]int, len(s.partialResults))

	for i := range s.partialResults {
		states[i] = s.partialResults[i].State
	}

	return WorstState(states...)
}
