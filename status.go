package check

import (
	"fmt"
	"strings"
)

const (
	// OK means everything is fine
	OK       = 0
	OKString = "OK"
	// Warning means there is a problem the admin should review
	Warning       = 1
	WarningString = "WARNING"
	// Critical means there is a problem that requires immediate action
	Critical       = 2
	CriticalString = "CRITICAL"
	// Unknown means the status can not be determined, probably due to an error or something missing
	Unknown       = 3
	UnknownString = "UNKNOWN"
)

// GetStatusText returns the string corresponding to a state
func GetStatusText(status int) (string, error) {
	switch status {
	case OK:
		return OKString, nil
	case Warning:
		return WarningString, nil
	case Critical:
		return CriticalString, nil
	case Unknown:
		return UnknownString, nil
	}

	return "", fmt.Errorf("no status text for status: %d", status)
}

// GetStatusInt returns a state corresponding to its
// common string representation
func GetStatusInt(status string) (int, error) {
	status = strings.ToUpper(status)

	switch status {
	case OKString, "0":
		return OK, nil
	case WarningString, "1":
		return Warning, nil
	case CriticalString, "2":
		return Critical, nil
	case UnknownString, "3":
		return Unknown, nil
	}

	return -1, fmt.Errorf("no integer for status: %s", status)
}
