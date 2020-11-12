package check

import (
	"errors"
	"strconv"
)

type nagiosPerfdataFloat struct {
	label string
	value float64
	uom	string
	warn rangeType
	crit rangeType
	min float64
	max float64
}

func (data nagiosPerfdataFloat) String() string {

	// Label
	label := formatLabel(&data.label)

	// Value
	value := strconv.FormatFloat(data.value, 'd', 2, 10)

	// UOM
	uom := data.uom

	// warn
	warn := formatRange(data.warn)

	// crit
	crit := formatRange(data.crit)

	// min + max
	min := strconv.FormatFloat(data.min, 'd', 2, 10)
	max := strconv.FormatFloat(data.max, 'd', 2, 10)

	return label + "=" + value + uom + ";" + warn + ";" + crit + ";" + min + ";" + max
}

func (perfdata nagiosPerfdataFloat)SanityCheck() error {
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
