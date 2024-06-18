package perfdata

import (
	"fmt"
	"math"
	"strings"

	"github.com/NETWAYS/go-check"
)

// Replace not allowed characters inside a label
var replacer = strings.NewReplacer("=", "_", "`", "_", "'", "_", "\"", "_")

type InfValueError struct {
}

func (i InfValueError) Error() string {
	return "Performance data value is infinite"
}

type NanValueError struct {
}

func (i NanValueError) Error() string {
	return "Performance data value is NaN (not a number)"
}

// formatNumeric returns a string representation of various possible numerics
//
// This supports most internal types of Go and all fmt.Stringer interfaces.
// Returns an eror in some known cases where the value of a data type does not
// represent a valid measurement, e.g INF for floats
// This error can probably ignored in most cases and the perfdata point omitted,
// but silently dropping the value and returning the empty strings seems like bad style
func formatNumeric(value interface{}) (string, error) {
	switch v := value.(type) {
	case float64:
		if math.IsInf(v, 0) {
			return "", InfValueError{}
		}

		if math.IsNaN(v) {
			return "", NanValueError{}
		}

		return check.FormatFloat(v), nil
	case float32:
		if math.IsInf(float64(v), 0) {
			return "", InfValueError{}
		}

		if math.IsNaN(float64(v)) {
			return "", NanValueError{}
		}

		return check.FormatFloat(float64(v)), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case fmt.Stringer, string:
		return fmt.Sprint(v), nil
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
// on errors (occurs with invalid data, the empty string is returned
func (p Perfdata) String() string {
	tmp, _ := p.ValidatedString()
	return tmp
}

// ValidatedString returns the proper format for the plugin output
// Returns an eror in some known cases where the value of a data type does not
// represent a valid measurement, see the explanation for "formatNumeric" for
// perfdata values.
func (p Perfdata) ValidatedString() (string, error) {
	var sb strings.Builder

	// Add quotes if string contains any whitespace
	if strings.ContainsAny(p.Label, "\t\n\f\r ") {
		sb.WriteString(`'` + replacer.Replace(p.Label) + `'` + "=")
	} else {
		sb.WriteString(replacer.Replace(p.Label) + "=")
	}

	pfVal, err := formatNumeric(p.Value)
	if err != nil {
		return "", err
	}

	sb.WriteString(pfVal)
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
			pfVal, err := formatNumeric(value)
			// Attention: we ignore limits if they are faulty
			if err == nil {
				sb.WriteString(pfVal)
			}
		}
	}

	return strings.TrimRight(sb.String(), ";"), nil
}
