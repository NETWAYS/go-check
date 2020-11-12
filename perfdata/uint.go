package check

import (
	"errors"
	"strconv"
)

type nagiosPerfdataUint struct {
	label string
	value uint64
	uom	string
	warn rangeType
	crit rangeType
	min uint64
	max uint64
}

func (data nagiosPerfdataUint) String() string {

	// Label
	label := formatLabel(&data.label)

	// Value
	value := strconv.FormatUint(data.value, 10)

	// UOM
	uom := data.uom

	// warn
	warn := formatRange(data.warn)

	// crit
	crit := formatRange(data.crit)

	// min + max
	min := strconv.FormatUint(data.min, 10)
	max := strconv.FormatUint(data.max, 10)

	return label + "=" + value + uom + ";" + warn + ";" + crit + ";" + min + ";" + max
}

func (perfdata nagiosPerfdataUint)SanityCheck() error {
	// Label
	err := sanityCheckLabel(&perfdata.label)

	// UOM
	err, uom := sanityCheckUom(&perfdata.uom)

	// value
	if uom == uomPercent {
		if perfdata.value < 0 || perfdata.value > 100 {
			return errors.New("Value not in percentage range")
		}
	}

	// warn
	err = sanityCheckRange(perfdata.warn)
	if err != nil {
		return err
	}

	// crit
	err = sanityCheckRange(perfdata.crit)
	if err != nil {
		return err
	}

	// min
	// max

	return nil
}
