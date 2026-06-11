package convert

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ByteIECUnits lists known units we can convert to based on uint64.
var ByteIECUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

// ByteSIUnits lists known units we can convert to based on uint64.
var ByteSIUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// IECBase is the exponential base for IEC units.
const IECBase = 1024

// SIBase is the exponential base for SI units.
const SIBase = 1000

// ParseBytes parses a strings and returns the number of bytes
func ParseBytes(value string) (uint64, error) {
	value = strings.TrimSpace(value)

	// Split number and unit by first non-numeric rune
	firstNonNumeric := func(c rune) bool { return !(c >= '0' && c <= '9' || c == '.') } //nolint: staticcheck

	i := strings.IndexFunc(value, firstNonNumeric)

	var unit string

	if i > 0 {
		unit = strings.TrimSpace(value[i:])
		value = value[0:i]
	}

	// Parse value to float64
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("provided value could not be parsed as float64: %s", value)
	}

	if number < 0 {
		return 0, fmt.Errorf("byte value cannot be negative: %s", value)
	}

	// Assume byte when no unit given
	if unit == "" {
		unit = "B"
	}

	// check for known units in ByteIECUnits
	for exponent, u := range ByteIECUnits {
		if u == unit {
			result := number * math.Pow(IECBase, float64(exponent))

			if result > math.MaxUint64 {
				return 0, fmt.Errorf("provided value could not be parsed as it overflows uint64: %s", value)
			}

			return uint64(result), nil
		}
	}

	// check for known units in ByteSIUnits
	for exponent, u := range ByteSIUnits {
		if u == unit {
			result := number * math.Pow(SIBase, float64(exponent))

			if result > math.MaxUint64 {
				return 0, fmt.Errorf("provided value could not be parsed as it overflows uint64: %s", value)
			}

			return uint64(result), nil
		}
	}

	return 0, fmt.Errorf("invalid unit: %s", unit)
}

// BytesSI returns the biggest sensible unit for the byte value with 2 decimal precision.
// When value is smaller than 2 render it with a lower scale.
func BytesSI(b uint64) string {
	value, unit := humanReadable(b, ByteSIUnits, SIBase)

	// Remove trailing zero decimals and any left over decimal dot
	s := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 2, 64), "0"), ".")

	return s + unit
}

// BytesIEC returns the biggest sensible unit for the byte value with 2 decimal precision.
// When value is smaller than 2 render it with a lower scale.
func BytesIEC(b uint64) string {
	value, unit := humanReadable(b, ByteIECUnits, IECBase)

	// Remove trailing zero decimals and any left over decimal dot
	s := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 2, 64), "0"), ".")

	return s + unit
}

// humanReadable searches for the closest feasible unit for displaying the byte value to a human.
//
// Meant as a universal function to be used by the implementations, with base and a list of unit names.
//
// A special behavior is that resulting values smaller than 2 are displayed with the lower exponent.
// If the input value is 0, humanReadable will always return "0B"
//
// Examples:
//
//	1073741824B -> 1000KB
//	2147483648B -> 2MB
//	0 -> 0B
func humanReadable(b uint64, units []string, base float64) (float64, string) {
	if b == 0 {
		return 0, "B"
	}

	exponent := math.Log(float64(b)) / math.Log(base)

	// Round to the unit scaled exponent
	unitExponent := math.Floor(exponent)

	// Ensure we only scale to the maximum known unit
	maxScale := float64(len(units) - 1)
	if unitExponent > maxScale {
		unitExponent = maxScale
	}

	value := math.Pow(base, exponent-unitExponent)

	// When resulting value is smaller than 2 calculate 1XXXM(i)B instead of 1.XXG(i)B
	if unitExponent > 0 && math.Round(value*base)/base < 2.0 {
		unitExponent--
		value = math.Pow(base, exponent-unitExponent)
	}

	return value, units[int(unitExponent)]
}
