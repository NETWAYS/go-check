package perfdata

import (
	"errors"
	"strconv"
)

type NagiosPerfdataFloat struct {
	Label string
	Value float64
	Uom	string
	Warn rangeType
	Crit rangeType
	Min float64
	Max float64
}

func (data NagiosPerfdataFloat) String() string {

	// Label
	label := formatLabel(&data.Label)

	// Value
	value := strconv.FormatFloat(data.Value, 'd', 2, 10)

	// UOM
	uom := data.Uom

	// Warn
	warn := formatRange(data.Warn)

	// Crit
	crit := formatRange(data.Crit)

	result := label + "=" + value + uom + ";" + warn + ";" + crit + ";"

	// Min + Max
	if data.Min != 0 {
		min := strconv.FormatFloat(data.Min, 'd', 2, 10)
		result += min + ";"
	} else {
		result += ";"
	}

	if data.Max != 0 {
		max := strconv.FormatFloat(data.Max, 'd', 2, 10)
		result += max
	}

	return result
}

func (perfdata NagiosPerfdataFloat)SanityCheck() error {
	// Label
	err := sanityCheckLabel(&perfdata.Label)

	// UOM
	err, Uom := sanityCheckUom(&perfdata.Uom)

	// Value
	if Uom == uomPercent {
		if perfdata.Value < 0 || perfdata.Value > 100 {
			return errors.New("Value not in percentage range")
		}
	}

	// Warn
	err = sanityCheckRange(perfdata.Warn)
	if err != nil {
		return err
	}

	// Crit
	err = sanityCheckRange(perfdata.Crit)
	if err != nil {
		return err
	}

	// Min
	// Max

	return nil
}
