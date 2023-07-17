package perfdata

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"strings"
)

// Replace not allowed characters inside a label
var replacer = strings.NewReplacer("=", "_", "`", "_", "'", "_", "\"", "_")

// formatNumeric returns a string representation of various possible numerics
//
// This supports most internal types of Go and all fmt.Stringer interfaces.
func formatNumeric(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return check.FormatFloat(v)
	case float32:
		return check.FormatFloat(float64(v))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case fmt.Stringer, string:
		return fmt.Sprint(v)
	default:
		panic(fmt.Sprintf("unsupported type for perfdata: %T", value))
	}
}

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
func (p Perfdata) String() string {
	var sb strings.Builder

	// Add quotes if string contains any whitespace
	if strings.ContainsAny(p.Label, "\t\n\f\r ") {
		sb.WriteString(`'` + replacer.Replace(p.Label) + `'` + "=")
	} else {
		sb.WriteString(replacer.Replace(p.Label) + "=")
	}

	sb.WriteString(formatNumeric(p.Value))
	sb.WriteString(p.Uom)

	// Thresholds
	for _, value := range []*check.Threshold{p.Warn, p.Crit} {
		sb.WriteString(";")

		if value != nil {
			sb.WriteString(value.String())
		}
	}

	// Limits
	for _, value := range []interface{}{p.Min, p.Max} {
		sb.WriteString(";")

		if value != nil {
			sb.WriteString(formatNumeric(value))
		}
	}

	return strings.TrimRight(sb.String(), ";")
}
