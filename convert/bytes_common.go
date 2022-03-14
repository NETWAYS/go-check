package convert

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type ByteAny interface {
	HumanReadable() string
	Bytes() uint64
	fmt.Stringer
}

// ExponentialFold defines the value to fold back to the previous exponent for better display of a small value.
const ExponentialFold = 2

// ParseBytes parses a strings and returns the proper ByteAny implementation based on the unit specified.
//
// When a plain value or B unit is used, ByteIEC is returned.
func ParseBytes(value string) (ByteAny, error) {
	value = strings.TrimSpace(value)

	// Split number and unit by first non-numeric rune
	firstNonNumeric := func(c rune) bool { return !(c >= '0' && c <= '9' || c == '.') }
	i := strings.IndexFunc(value, firstNonNumeric)

	var unit string

	if i > 0 {
		unit = strings.TrimSpace(value[i:])
		value = value[0:i]
	}

	// Parse value to float64
	number, err := strconv.ParseFloat(value, 64) // nolint:gomnd
	if err != nil {
		return nil, fmt.Errorf("provided value could not be parsed as float64: %s", value)
	}

	// Assume byte when no unit given
	if unit == "" {
		unit = "B"
	}

	// check for known units in ByteIECUnits
	for exponent, u := range ByteIECUnits {
		if u == unit {
			// convert to bytes and return type
			return BytesIEC(number * math.Pow(IECBase, float64(exponent))), nil
		}
	}

	// check for known units in ByteSIUnits
	for exponent, u := range ByteSIUnits {
		if u == unit {
			// convert to bytes and return type
			return BytesSI(number * math.Pow(SIBase, float64(exponent))), nil
		}
	}

	return nil, fmt.Errorf("invalid unit: %s", unit)
}

// humanReadable searches for the closest feasible unit for displaying the byte value to a human.
//
// Meant as a universal function to be used by the implementations, with base and a list of unit names.
//
// A special behavior is that resulting values smaller than 2 are displayed with the lower exponent.
// If the input value is 0, humanReadable will always return "0MB"
//
// Examples:
//  1073741824B -> 1000KB
//  2147483648B -> 2MB
//  0 -> 0MB
//
func humanReadable(b uint64, units []string, base float64) (float64, string) {
	if b == 0 {
		return 0, "MB"
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
	if math.Round(value*base)/base < ExponentialFold {
		unitExponent--
		value = math.Pow(base, exponent-unitExponent)
	}

	return value, units[int(unitExponent)]
}
