package check

import (
	"errors"
	"strconv"
)

type nagiosPerfdataInt struct {
	label string
	value int64
	uom	string
	warn rangeType
	crit rangeType
	min int64
	max int64
}

func (data nagiosPerfdataInt) String() string {

	// Label
	label := formatLabel(&data.label)

	// Value
	value := strconv.FormatInt(data.value, 10)

	// UOM
	uom := data.uom

	// warn
	warn := formatRange(data.warn)

	// crit
	crit := formatRange(data.crit)

	// min + max
	min := strconv.FormatInt(data.min, 10)
	max := strconv.FormatInt(data.max, 10)

	return label + "=" + value + uom + ";" + warn + ";" + crit + ";" + min + ";" + max
}

func (perfdata nagiosPerfdataInt)SanityCheck() error {
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
