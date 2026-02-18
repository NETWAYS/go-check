package result

import "github.com/NETWAYS/go-check"

// WorstState determines the worst state from a list of states
//
// Helps combining an overall states, only based on a
// few numbers for various checks.
//
// Order of preference: Critical, Unknown, Warning, Ok
func WorstState(states ...check.Status) check.Status {
	var overall check.Status

	overall = check.OK
	// nolint: gocritic
	for _, state := range states {
		if state == check.Critical {
			overall = check.Critical
		} else if state == check.Unknown {
			if overall != check.Critical {
				overall = check.Unknown
			}
		} else if state > overall {
			overall = state
		}
	}

	return overall
}
