// result tries to
package result

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/perfdata"
	"strings"
)

type Overall struct {
	OKs       int
	Warnings  int
	Criticals int
	Unknowns  int
	Summary   string
	Outputs   []string // Deprecate this in a future version
	subchecks []Subcheck
}

type Subcheck struct {
	State int
	Output string
	Perfdata perfdata.PerfdataList
	subchecks []Subcheck
}

func (s *Subcheck) String() string {
	return fmt.Sprintf("[%s] %s|%s", check.StatusText(s.State), s.Output, s.Perfdata.String())
}

// Deprecate this in a future version
func (o *Overall) AddOK(output string) {
	o.Add(check.OK, output)
}

// Deprecate this in a future version
func (o *Overall) AddWarning(output string) {
	o.Add(check.Warning, output)
}

// Deprecate this in a future version
func (o *Overall) AddCritical(output string) {
	o.Add(check.Critical, output)
}

// Deprecate this in a future version
func (o *Overall) AddUnknown(output string) {
	o.Add(check.Unknown, output)
}

func (o *Overall) Add(state int, output string) {
	switch state {
	case check.OK:
		o.OKs++
	case check.Warning:
		o.Warnings++
	case check.Critical:
		o.Criticals++
	default:
		o.Unknowns++
	}

	o.Outputs = append(o.Outputs, fmt.Sprintf("[%s] %s", check.StatusText(state), output))
}

func (o* Overall) AddSubcheck(subcheck Subcheck) {
	o.subchecks = append(o.subchecks, subcheck)
}

func (o* Subcheck) AddSubcheck(subcheck Subcheck) {
	o.subchecks = append(o.subchecks, subcheck)
}

func (o *Overall) GetStatus() int {
	if o.Criticals > 0 {
		return check.Critical
	} else if o.Unknowns > 0 {
		return check.Unknown
	} else if o.Warnings > 0 {
		return check.Warning
	} else if o.OKs > 0 {
		return check.OK
	} else {
		return check.Unknown
	}
}

func (o *Overall) GetSummary() string {
	if o.Summary == "" {
		if o.Criticals > 0 {
			o.Summary += fmt.Sprintf("critical=%d ", o.Criticals)
		}

		if o.Unknowns > 0 {
			o.Summary += fmt.Sprintf("unknown=%d ", o.Unknowns)
		}

		if o.Warnings > 0 {
			o.Summary += fmt.Sprintf("warning=%d ", o.Warnings)
		}

		if o.OKs > 0 {
			o.Summary += fmt.Sprintf("ok=%d ", o.OKs)
		}

		if o.Summary == "" && len(o.subchecks) == 0 {
			o.Summary = "No status information"
		} else {
			criticals := 0
			warnings := 0
			oks := 0
			unknowns := 0
			for _, sc := range o.subchecks {
				if sc.State == check.Critical {
					criticals ++
				} else if sc.State == check.Warning {
					warnings ++
				} else if sc.State == check.Unknown {
					unknowns ++
				} else if sc.State == check.OK {
					oks ++
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
			o.Summary = "states: " + strings.TrimSpace(o.Summary)
		}
	}

	return o.Summary
}

func (o *Overall) GetOutput() string {
	output := o.GetSummary() + "\n"

	for _, extra := range o.Outputs {
		output += extra + "\n"
	}

	if o.subchecks != nil {
		for _, s := range o.subchecks {
			output += s.getOutput(0)
		}
	}

	return output
}

func (s *Subcheck) getOutput(indent_level int) string {
	var output string
	prefix := strings.Repeat("  ", indent_level)
	output += prefix + "|- " + s.String() + "\n"

	if s.subchecks != nil {
		for _, ss := range s.subchecks {
			output += ss.getOutput(indent_level + 1)
		}
	}

	return output
}
