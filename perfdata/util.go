package perfdata

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"regexp"
	"strings"
)

// Lists all allowed characters inside a label, so we can replace any non-matching
var validInLabelRe = regexp.MustCompile(`[^a-zA-Z0-9 _\-+:/.;]+`)

// FormatNumeric returns a string representation of various possible numerics
//
// This supports most internal types of Go and all fmt.Stringer interfaces.
func FormatNumeric(value interface{}) string {
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
