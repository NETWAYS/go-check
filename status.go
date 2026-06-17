package check

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// OKString means everything is fine
	OKString = "OK"
	// WarningString means there is a problem the admin should review
	WarningString = "WARNING"
	// CriticalString means there is a problem that requires immediate action
	CriticalString = "CRITICAL"
	// UnknownString means the status can not be determined, probably due to an error or something missing
	UnknownString = "UNKNOWN"
)

type Status int

const (
	OK Status = iota
	Warning
	Critical
	Unknown
)

// NewStatus returns a state corresponding to its
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

	return Unknown, fmt.Errorf("%d is not a valid state", status)
}

// NewStatusFromString returns a state corresponding to its
// common string representation
func NewStatusFromString(status string) (Status, error) {
	s := strings.ToUpper(status)
	switch s {
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
	default:
		return UnknownString
	}
}

// WorstState determines the worst state from a list of states
//
// Helps combining an overall states, only based on a
// few numbers for various checks.
//
// Order of preference: Critical, Unknown, Warning, Ok
func WorstState(states ...Status) Status {
	if len(states) < 1 {
		return Unknown
	}

	overall := OK

	// nolint: gocritic
	for _, state := range states {
		if Compare(overall, state) > 0 {
			overall = state
		}
	}

	if overall < OK || overall > Unknown {
		overall = Unknown
	}

	return overall
}

// Compare compares two Status types
// if the left one (a) is worse than the right one (b), the result is < 0
// if they are equal, the result is 0
// if the right one (b) is worse than the left one (a), the result is > 0
func Compare(a Status, b Status) int {
	switch a {
	case OK:
		switch b {
		case OK:
			return 0
		case Warning, Unknown, Critical:
			return 1
		}
	case Warning:
		switch b {
		case OK:
			return -1
		case Warning:
			return 0
		case Unknown, Critical:
			return 1
		}
	case Unknown:
		switch b {
		case OK, Warning:
			return -1
		case Unknown:
			return 0
		case Critical:
			return 1
		}
	case Critical:
		switch b {
		case OK, Warning, Unknown:
			return -1
		case Critical:
			return 0
		}
	}

	// should not be possible to land here
	return 0
}
