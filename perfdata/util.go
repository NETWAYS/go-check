package perfdata

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidUom lists all valid units allowed by spec and Icinga 2
//
// Also see:
//   - https://www.monitoring-plugins.org/doc/guidelines.html#AEN201
//   - https://github.com/Icinga/icinga2/blob/master/lib/base/perfdatavalue.cpp
const ValidUom = "us|ms|s|tb|gb|mb|kb|b|%|c"

// Lists all allowed characters inside a label, so we can replace any non-matching
var validInLabelRe = regexp.MustCompile(`[^a-zA-Z0-9 _\-+:/.]+`)

var validUomSlice = strings.Split(ValidUom, "|")

// FormatNumeric returns a string representation of various possible numerics
//
// This supports most internal types of Go and all fmt.Stringer interfaces.
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

// FormatLabel returns a sane perfdata label
//
// All groups of invalid characters will be replaced by a single underscore.
func FormatLabel(label string) string {
	// Replace invalid character groups by an underscore
	label = validInLabelRe.ReplaceAllString(label, "_")

	if strings.ContainsAny(label, " ") {
		return `'` + label + `'`
	}

	return label
}

// IsValidUom compares the given unit against ValidUom
func IsValidUom(uom string) bool {
	for _, k := range validUomSlice {
		if k == uom {
			return true
		}
	}

	return false
}
