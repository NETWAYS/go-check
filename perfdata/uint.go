package perfdata

import (
	"strconv"
)

type NagiosPerfdataUint struct {
	Label string
	Value uint64
	Uom   string
	Warn  Threshold
	Crit  Threshold
	Min   uint64
	Max   uint64
}

func (data NagiosPerfdataUint) String() string {

	// Label
	label := formatLabel(&data.Label)

	// Value
	value := strconv.FormatUint(data.Value, 10)

	// UOM
	uom := data.Uom

	// warn
	warn := formatRange(data.Warn)

	// crit
	crit := formatRange(data.Crit)

	result := label + "=" + value + uom + ";" + warn + ";" + crit + ";"

	// Min + Max
	if data.Min != 0 {
		Min := strconv.FormatUint(data.Min, 10)
		result += Min + ";"
	} else {
		result += ";"
	}

	if data.Max != 0 {
		Max := strconv.FormatUint(data.Max, 10)
		result += Max
	}

	return result
}

/*
func (perfdata NagiosPerfdataUint)SanityCheck() error {
	// Label
	err := sanityCheckLabel(&perfdata.Label)

	// UOM
	err, uom := sanityCheckUom(&perfdata.Uom)

	// value
	if uom == uomPercent {
		if perfdata.Value < 0 || perfdata.Value > 100 {
			return errors.New("Value not in percentage range")
		}
	}

	// warn
	err = sanityCheckRange(perfdata.Warn)
	if err != nil {
		return err
	}

	// crit
	err = sanityCheckRange(perfdata.Crit)
	if err != nil {
		return err
	}

	// Min
	// Max

	return nil
}
*/
