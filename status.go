package check

import (
	"fmt"
	"strings"
)

type Status int

const (
	OK Status = iota
	Warning
	Critical
	Unknown
)

const (
	OKString       = "OK"
	WarningString  = "WARNING"
	CriticalString = "CRITICAL"
	UnknownString  = "UNKNOWN"
)

func (s Status) String() string {
	switch s {
	case OK:
		return OKString
	case Warning:
		return WarningString
	case Critical:
		return CriticalString
	case Unknown:
		return UnknownString
	}

	return UnknownString
}

// NewStatusFromInt returns a state corresponding to its interger
func NewStatusFromInt(status int) (Status, error) {
	switch status {
	case 0:
		return OK, nil
	case 1:
		return Warning, nil
	case 2:
		return Critical, nil
	case 3:
		return Unknown, nil
	default:
		return Unknown, fmt.Errorf("unable to create status of type %d", status)
	}
}

// NewStatusFromString returns a state corresponding to its
// common string representation
func NewStatusFromString(status string) (Status, error) {
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
	default:
		return Unknown, fmt.Errorf("unable to create status of type %s", status)
	}
}
