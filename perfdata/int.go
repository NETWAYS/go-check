package perfdata

import (
	"errors"
	"strconv"
)

type NagiosPerfdataInt struct {
	Label string
	Value int64
	Uom	string
	Warn rangeType
	Crit rangeType
	Min int64
	Max int64
}

func (data NagiosPerfdataInt) String() string {

	// Label
	label := formatLabel(&data.Label)

	// Value
	value := strconv.FormatInt(data.Value, 10)

	// UOM
	uom := data.Uom

	// Warn
	warn := formatRange(data.Warn)

	// Crit
	crit := formatRange(data.Crit)

	result := label + "=" + value + uom + ";" + warn + ";" + crit + ";"

	// Min + Max
	if data.Min != 0 {
		Min := strconv.FormatInt(data.Min, 10)
		result += Min + ";"
	} else {
		result += ";"
	}

	if data.Max != 0 {
		Max := strconv.FormatInt(data.Max, 10)
		result += Max
	}

	return result
}

func (perfdata NagiosPerfdataInt)SanityCheck() error {
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
