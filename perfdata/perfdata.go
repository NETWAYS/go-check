package perfdata

import (
	"errors"
	"regexp"
	"strings"
)

const (
	uomNone = iota
	uomByte
	uomPercent
	uomTime
	uomError
)

type uomType uint

type perfdata interface {
	SanityCheck() error
	String() string
}

func formatLabel(label *string) string {
	// Label
	var result string
	if strings.ContainsAny(*label, " \t") {
		result = "'" + *label + "'"
	} else {
		result = *label
	}
	return result
}

// --- Try to specify Collection format
// PERFDATA := (Label=ValueUnitofmeasurement;Warn;Crit;Min;Max)+
// Label := [!'=]+ // All characters except ' and = , but I also don't want \n
// Value := INT
// Unitofmeasurement := EMPTY|SECOND|PERCENT|BYTES
// Warn := EMPTY|RANGE
// Crit := EMPTY|RANGE
// Min := EMPTY|INT
// Max := EMPTY|INT
// RANGE := [@]StartEnd // start must be <= end
// START := EMPTY|INT:|~:
// END := EMPTY|INT|~
// INT = -*[0-9]+ // Integer, may be negative
// EMPTY = "" // Empty String
// SECOND := s|us|ms
// BYTES := B|KB|MB|TB
// PERCENT := %
// Note on range: The specification above does not reflect the logic of the format
// if start == 0 start and the following : may be omitted, the end may also be omitted if the end is infinity
// the tilde character specifies negative infinity
// if the range starts with @ the alert is inside the range (inclusive endpoints), otherwise Inside (inclusive endpoints)
// ---
// Check input according to https://www.monitoring-plugins.org/doc/guidelines.html#AEN201
// or https://nagios-plugins.org/doc/guidelines.html

func sanityCheckLabel(label *string) error {
	// label
	if strings.ContainsAny(*label, "'=\n") {
		// Restricted character in label!
		return errors.New("Illegal character in perfdata label: " + *label)
	}
	if len(*label) == 0 {
		return errors.New("No label given")
	}
	return nil
}

func sanityCheckUom(uom *string) (error, uomType) {
	if *uom == "" {
		return nil, uomNone
	}
	if *uom == "%" {
		return nil, uomPercent
	}

	match, err := regexp.MatchString("^(s|us|ms)$", *uom)
	if err != nil {
		return err, uomError
	}
	if !match {
	} else {
		return nil, uomTime
	}

	match, err = regexp.MatchString("^(B|KB|MB|TB)$", *uom)
	if err != nil {
		return err, uomByte
	}
	if !match {
		return errors.New("No matching unit of Measurement"), uomError
	} else {
		return nil, uomError
	}
}

/*
func sanityCheckRange(rangeValue Threshold) error {
	if rangeValue.Lower == "~" {
		// this is ok
		if rangeValue.Upper != "~" {
			// Since start <= end this can only be wrong
			return errors.New("Threshold Error: Start > End! This is wrong")
		} else {
			// This is valid, although useless
			// Warning: Mathematicans might disagree with that
			return nil
		}
	} else if rangeValue.Lower == "" {
		// This is equivalent to Lower = 0
		// So, Upper must be >0
		// this be infinity or a number value
		if rangeValue.Upper == "" {
			// infty, this is fine
			return nil
		}
		if num, err := strconv.ParseInt(rangeValue.Upper, 10, 64); err == nil {
			if num < 0 {
				return errors.New("Threshold Error: End < Start")
			}
		}
		if num, err := strconv.ParseFloat(rangeValue.Upper, 64); err == nil {
			if num < 0 {
				return errors.New("Threshold Error: End < Start")
			}
		} else {
			return  errors.New("Threshold Error: Could not parse upper Bound")
		}
	}

	// At this point there has to a number in Lower
	if lower, err := strconv.ParseInt(rangeValue.Lower, 10, 64); err == nil {
		if upper, err := strconv.ParseInt(rangeValue.Upper, 10, 64); err == nil {
			if upper < lower  {
				return errors.New("Threshold Error: End < Start")
			} else {
				return nil
			}
		}
		if upper, err := strconv.ParseFloat(rangeValue.Upper, 64); err == nil {
			if upper < float64(lower) {
				return errors.New("Threshold Error: End < Start")
			} else {
				return nil
			}
		} else {
			return  errors.New("Threshold Error: Could not parse upper Bound")
		}
	}
	if lower, err := strconv.ParseFloat(rangeValue.Lower, 64); err == nil {
		if upper, err := strconv.ParseInt(rangeValue.Upper, 10, 64); err == nil {
			if float64(upper) < lower  {
				return errors.New("Threshold Error: End < Start")
			} else {
				return nil
			}
		}
		if upper, err := strconv.ParseFloat(rangeValue.Upper, 64); err == nil {
			if upper < float64(lower) {
				return errors.New("Threshold Error: End < Start")
			} else {
				return nil
			}
		} else {
			return  errors.New("Threshold Error: Could not parse upper Bound")
		}
	} else {
		return  errors.New("Threshold Error: Could not parse lower Bound")
	}
}
*/
