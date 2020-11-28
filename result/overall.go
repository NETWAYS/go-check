// result tries to
package result

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"strings"
)

type Overall struct {
	OKs       int
	Warnings  int
	Criticals int
	Unknowns  int
	Summary   string
	Outputs   []string
}

func (o *Overall) AddOK(output string) {
	o.Add(check.OK, output)
}

func (o *Overall) AddWarning(output string) {
	o.Add(check.Warning, output)
}

func (o *Overall) AddCritical(output string) {
	o.Add(check.Critical, output)
}

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
		if o.Summary == "" {
			o.Summary = "No status information"
		} else {
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

	return output
}
