package perfdata

import (
	"fmt"
	"regexp"
	"strings"
	"github.com/NETWAYS/go-check"
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

// Lists all allowed characters inside a label, so we can replace any non-matching
var validInLabelRe = regexp.MustCompile(`[^a-zA-Z0-9 _\-+:/.]+`)

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

func FormatNumeric(value interface{}) string {
	switch value.(type) {
	case float64, float32:
		return fmt.Sprintf("%g", value)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
		return fmt.Sprintf("%d", value)
	case fmt.Stringer, string:
		return fmt.Sprintf("%s", value)
	default:
		panic(fmt.Sprintf("unsupported type for perfdata: %T", value))
	}
}

func FormatLabel(label string) string {
	// Replace invalid character groups by an underscore
	label = validInLabelRe.ReplaceAllString(label, "_")

	if strings.ContainsAny(label, " ") {
		return fmt.Sprintf(`'%s'`, label)
	}

	return label
}
