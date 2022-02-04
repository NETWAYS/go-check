package perfdata

import (
	"github.com/NETWAYS/go-check"
	"strings"
)

// Perfdata represents all properties of performance data for Icinga
//
// Implements fmt.Stringer to return the plaintext format for a plugin output.
//
// For examples of Uom see:
//
// https://www.monitoring-plugins.org/doc/guidelines.html#AEN201
//
// https://github.com/Icinga/icinga2/blob/master/lib/base/perfdatavalue.cpp
//
// https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/#unit-of-measurement-uom
type Perfdata struct {
	Label string
	Value interface{}
	// Uom is the unit-of-measurement, see links above for details.
	Uom  string
	Warn *check.Threshold
	Crit *check.Threshold
	Min  interface{}
	Max  interface{}
}

// String returns the proper format for the plugin output
func (p Perfdata) String() (s string) {
	s = FormatLabel(p.Label) + "="

	s += FormatNumeric(p.Value)
	s += p.Uom

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
