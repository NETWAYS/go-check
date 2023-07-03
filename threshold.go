package check

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Defining a threshold for any numeric value
//
// Format: [@]start:end
//
// Threshold  Generate an alert if x...
// 10         < 0 or > 10, (outside the range of {0 .. 10})
// 10:        < 10, (outside {10 .. ∞})
// ~:10       > 10, (outside the range of {-∞ .. 10})
// 10:20      < 10 or > 20, (outside the range of {10 .. 20})
// @10:20     ≥ 10 and ≤ 20, (inside the range of {10 .. 20})
//
// Reference: https://www.monitoring-plugins.org/doc/guidelines.html#THRESHOLDFORMAT
type Threshold struct {
	Inside bool
	Lower  float64
	Upper  float64
}

var (
	thresholdNumberRe = regexp.MustCompile(`(-?\d+(?:\.\d+)?|~)`)
	thresholdRe       = regexp.MustCompile(fmt.Sprintf(`^(@)?(?:%s:)?(?:%s)?$`,
		thresholdNumberRe.String(), thresholdNumberRe.String()))
	PosInf = math.Inf(1)
	NegInf = math.Inf(-1)
)

// Parse a Threshold from a string.
//
// See Threshold for details.
func ParseThreshold(spec string) (t *Threshold, err error) {
	t = &Threshold{}

	parts := thresholdRe.FindStringSubmatch(spec)
	if spec == "" || len(parts) == 0 {
		err = fmt.Errorf("could not parse threshold: %s", spec)
		return
	}

	// @ at the beginning
	if parts[1] != "" {
		t.Inside = true
	}

	var v float64

	// Lower bound
	if parts[2] == "~" {
		t.Lower = NegInf
	} else if parts[2] != "" {
		v, err = strconv.ParseFloat(parts[2], 64)
		if err != nil {
			err = fmt.Errorf("can not parse lower bound '%s': %w", parts[2], err)
			return
		}

		t.Lower = v
	}

	// Upper bound
	if parts[3] == "~" || (parts[3] == "" && parts[2] != "") {
		t.Upper = PosInf
	} else if parts[3] != "" {
		v, err = strconv.ParseFloat(parts[3], 64)
		if err != nil {
			err = fmt.Errorf("can not parse upper bound '%s': %w", parts[3], err)
			return
		}

		t.Upper = v
	}

	return
}

// String returns the plain representation of the Threshold
func (t Threshold) String() (s string) {
	s = BoundaryToString(t.Upper)

	// remove upper ~, which is the default
	if s == "~" {
		s = ""
	}

	if t.Lower != 0 {
		s = BoundaryToString(t.Lower) + ":" + s
	}

	if t.Inside {
		s = "@" + s
	}

	return
}

// Compares a value against the threshold, and returns true if the value violates the threshold.
func (t Threshold) DoesViolate(value float64) bool {
	if t.Inside {
		return value >= t.Lower && value <= t.Upper
	} else {
		return value < t.Lower || value > t.Upper
	}
}

// BoundaryToString returns the string representation of a Threshold boundary.
func BoundaryToString(value float64) (s string) {
	s = FormatFloat(value)

	// In the threshold context, the sign derives from lower and upper bound, we only need the ~ notation
	if s == "+Inf" || s == "-Inf" {
		s = "~"
	}

	return
}

// FormatFloat returns a string representation of floats, avoiding scientific notation and removes trailing zeros.
func FormatFloat(value float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", value), "0"), ".") // remove trailing 0 and trailing dot
}
