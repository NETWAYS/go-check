package check

import (
	"errors"
	"fmt"
)

const (
	// Invalid is not a valid status
	InvalidString = "Invalid"
	// OK means everything is fine
	OKString = "OK"
	// Warning means there is a problem the admin should review
	WarningString = "WARNING"
	// Critical means there is a problem that requires immediate action
	CriticalString = "CRITICAL"
	// Unknown means the status can not be determined, probably due to an error or something missing
	UnknownString = "UNKNOWN"
)

type Status int

const (
	Invalid Status = iota - 1
	OK
	Warning
	Critical
	Unknown
)

// NewStatusFromString returns a state corresponding to its
// common string representation
func NewStatus(status int) (Status, error) {
	switch status {
	case 0:
		return OK, nil
	case 1:
		return Warning, nil
	case 2:
		return Critical, nil
	case 3:
		return Unknown, nil
	}

	return Invalid, fmt.Errorf("%d is not a valid state", status)
}

// NewStatusFromString returns a state corresponding to its
// common string representation
func NewStatusFromString(status string) (Status, error) {
	switch status {
	case OKString:
		return OK, nil
	case WarningString:
		return Warning, nil
	case CriticalString:
		return Critical, nil
	case UnknownString:
		return Unknown, nil
	}

	return Invalid, errors.New(status + " is not a valid state")
}

// String returns the string corresponding to a state
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

	return InvalidString
}
