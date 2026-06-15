package result

import "github.com/NETWAYS/go-check"

// WorstState determines the worst state from a list of states
//
// Helps combining an overall states, only based on a
// few numbers for various checks.
//
// Order of preference: Critical, Unknown, Warning, Ok
func WorstState(states ...check.Status) check.Status {
	if len(states) < 1 {
		return check.Unknown
	}

	overall := check.OK

	// nolint: gocritic
	for _, state := range states {
		if check.Compare(overall, state) > 0 {
			overall = state
		}
	}

	if overall < check.OK || overall > check.Unknown {
		overall = check.Unknown
	}

	return overall
}
