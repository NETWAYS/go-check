package perfdata

import (
	"github.com/NETWAYS/go-check"
	"strings"
)

// Perfdata represents all properties of performance data for Icinga
//
// Implements fmt.Stringer to return the plaintext format for a plugin output.
//
// Also see https://www.monitoring-plugins.org/doc/guidelines.html#AEN201
type Perfdata struct {
	Label string
	Value interface{}
	Uom   string
	Warn  *check.Threshold
	Crit  *check.Threshold
	Min   interface{}
	Max   interface{}
}

// String returns the proper format for the plugin output
func (p Perfdata) String() (s string) {
	s = FormatLabel(p.Label) + "="

	// Value
	s += FormatNumeric(p.Value)

	if IsValidUom(p.Uom) {
		s += p.Uom
	}

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

	return
}
