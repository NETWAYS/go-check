package perfdata

import (
	"github.com/NETWAYS/go-check"
	"strings"
)

type Perfdata struct {
	Label string
	Value interface{}
	Uom   string
	Warn  *check.Threshold
	Crit  *check.Threshold
	Min   interface{}
	Max   interface{}
}

func (p Perfdata) String() (s string) {
	s = FormatLabel(p.Label) + "="

	// Value
	s += FormatNumeric(p.Value)
	s += p.Uom // TODO: typing and nil check?

	// Thresholds
	for _, value := range []*check.Threshold{p.Warn, p.Crit} {
		s += ";"
		if value != nil {
			s += value.String()
		}
	}

	// Limits
	for _, value := range []interface{}{p.Min, p.Max} {
		s += ";"
		if value != nil {
			s += FormatNumeric(value)
		}
	}

	// Remove trailing semicolons
	s = strings.TrimRight(s, ";")

	return
}
