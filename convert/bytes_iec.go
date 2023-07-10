package convert

import (
	"strconv"
	"strings"
)

// ByteIECUnits lists known units we can convert to based on uint64.
//
// See https://en.wikipedia.org/wiki/Byte#Unit_symbol
var ByteIECUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

// BytesIEC is the IEC (1024) unit implementation of byte conversion.
type BytesIEC uint64

// IECBase is the exponential base for IEC units.
const IECBase = 1024

// HumanReadable returns the biggest sensible unit for the byte value with 2 decimal precision.
//
// When value is smaller than 2 render it with a lower scale.
func (b BytesIEC) HumanReadable() string {
	value, unit := humanReadable(uint64(b), ByteIECUnits, IECBase)

	// Remove trailing zero decimals and any left over decimal dot
	s := strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 2, 64), "0"), ".")

	return s + unit
}

func (b BytesIEC) String() string {
	return b.HumanReadable()
}

// Bytes returns the value as uint64.
func (b BytesIEC) Bytes() uint64 {
	return uint64(b)
}
