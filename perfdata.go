package check

import (
	"strings"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// return a correctly formatted perfdata string
func FormatDataPoint(label string, value string, unitOfMeasurement string, warn string, crit string, min string, max string) (string, error) {
	// Uncomment this for validating the format
	// if debug {
	// 		SanityCheckPerfData(label, value, unitOfMeasurement, warn, crit, min, max)
	// }
	if strings.ContainsAny(label, " \t") {
		label = "'" + label + "'"
	}

	result := fmt.Sprintf("%v=%v%v;%v;%v;%v;%v", label, value, unitOfMeasurement, warn, crit, min, max)
	return result, nil
}

// --- Try to specify Perfdata format
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
// if the range starts with @ the alert is inside the range (inclusive endpoints), otherwise outside (inclusive endpoints)
// ---
// Check input according to https://www.monitoring-plugins.org/doc/guidelines.html#AEN201
func SanityCheckPerfData(label string, value string, unitOfMeasurement string, warn string, crit string, min string, max string) error {

	// label
	if strings.ContainsAny(label, "'=\n") {
		// Restricted character in label!
		return errors.New("Illegal character in perfdata label: " + label)
	}
	if len(label) == 0 {
		return errors.New("No label given")
	}

	// value
	match, err := regexp.MatchString("-*[0-9]+", value)
	if err != nil {
		return err
	}
	if !match && value != "U" {
		return errors.New("Non digit symbols in value")
	}

	// UOM
	if unitOfMeasurement  != "" {
		if unitOfMeasurement == "%" {
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			if i < 0 || i > 100 {
				return errors.New("value not in percent range")
			}
		}

		match, err := regexp.MatchString("^(s|us|ms|B|KB|MB|TB)$", unitOfMeasurement)
		if err != nil {
			return err
		}

		if !match {
			return errors.New("No matching unit of Measurement")
		}
	}

	// warn
	err = validateRange(warn)
	if err != nil {
		return err
	}

	// crit
	err = validateRange(crit)
	if err != nil {
		return err
	}

	// min
	match, err = regexp.MatchString("|(-*[0-9]+)", min)
	if err != nil {
		return err
	}

	// max
	match, err = regexp.MatchString("|(-*[0-9]+)", max)
	if err != nil {
		return err
	}

	return nil
}

func validateRange(rangeString string) error {
	if rangeString[0] == '@' {
		// Don't care if alert is inside range or not
		rangeString = rangeString[1:]
	}

	if !strings.Contains(rangeString, ":") {
		// no colon in string, so start was omitted
		match, err := regexp.MatchString("-*[0-9]+", rangeString)
		if err != nil {
			return err
		}
		if !match {
			return errors.New("Range value contains wrong characters: " + rangeString)
		}
	}

	match, err := regexp.MatchString("(~|[0-9]+):(~|[0-9]+)", rangeString)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("Range format was invalid: " + rangeString)
	}

	parts := strings.Split(rangeString, ":")

	// start
	if parts[0] != "~" {
		match, err := regexp.MatchString("-*[0-9]+", parts[0])
		if err != nil {
			return err
		}
		if !match {
			return errors.New("Range start does not match format specification: " + parts[0])
		}
	}
	// end
	if parts[1] != "~" {
		match, err := regexp.MatchString("-*[0-9]+", parts[1])
		if err != nil {
			return err
		}
		if !match {
			return errors.New("Range start does not match format specification: " + parts[1])
		}
	}

	return nil
}
