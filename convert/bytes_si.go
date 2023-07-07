package convert

import (
	"strconv"
	"strings"
)

// ByteSIUnits lists known units we can convert to based on uint64.
//
// See https://en.wikipedia.org/wiki/Byte#Unit_symbol
var ByteSIUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// BytesSI is the SI (1000) unit implementation of byte conversion.
type BytesSI uint64

// SIBase is the exponential base for SI units.
const SIBase = 1000

// HumanReadable returns the biggest sensible unit for the byte value with 2 decimal precision.
//
// When value is smaller than 2 render it with a lower scale.
func (b BytesSI) HumanReadable() string {
	value, unit := humanReadable(uint64(b), ByteSIUnits, SIBase)

	// Remove trailing zero decimals and any left over decimal dot
	s := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 2, 64), "0"), ".")

	return s + unit
}

func (b BytesSI) String() string {
	return b.HumanReadable()
}

// Bytes returns the value as uint64.
func (b BytesSI) Bytes() uint64 {
	return uint64(b)
}
