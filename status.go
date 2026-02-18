package check

import (
	"errors"
)

const (
	OKString       = "OK"
	WarningString  = "WARNING"
	CriticalString = "CRITICAL"
	UnknownString  = "UNKNOWN"
)

type Status int

const (
	// OK means everything is fine
	OK = iota
	// Warning means there is a problem the admin should review
	Warning
	// Critical means there is a problem that requires immediate action
	Critical
	// Unknown means the status can not be determined, probably due to an error or something missing
	Unknown
)

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

	return Unknown, errors.New(status + " is not a valid state")
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

	return UnknownString
}
