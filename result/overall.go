// result tries to
package result

import (
	"fmt"
	check "go-check"
	"go-check/perfdata"
	"strings"
)

type Overall struct {
	OKs       int
	Warnings  int
	Criticals int
	Unknowns  int
	Summary   string
	Outputs   []string
	PerfdataUint []perfdata.NagiosPerfdataUint
	PerfdataInt []perfdata.NagiosPerfdataInt
	PerfdataFloat []perfdata.NagiosPerfdataFloat
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

	output += o.formatPerfdata()

	return output
}

func (o *Overall) formatPerfdata() string {
	if (len(o.PerfdataInt) == 0) && (len(o.PerfdataUint) == 0) && (len(o.PerfdataFloat) == 0) {
		return ""
	}
	result := "|"
	for _, pdInt := range o.PerfdataInt {
		result += pdInt.String() + " "
	}
	for _, pdUint := range o.PerfdataUint {
		result += pdUint.String() + " "
	}
	for _, pdFloat := range o.PerfdataFloat {
		result += pdFloat.String() + " "
	}
	return result
}

func (o *Overall) AddNagiosPerfdataInt(data perfdata.NagiosPerfdataInt) {
	o.PerfdataInt = append(o.PerfdataInt, data)
}
func (o *Overall) AddNagiosPerfdataUint(data perfdata.NagiosPerfdataUint) {
	o.PerfdataUint = append(o.PerfdataUint, data)
}
func (o *Overall) AddNagiosPerfdataFloat(data perfdata.NagiosPerfdataFloat) {
	o.PerfdataFloat = append(o.PerfdataFloat, data)
}
