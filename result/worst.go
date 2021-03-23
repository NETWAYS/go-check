package result

import "github.com/NETWAYS/go-check"

// Determines the worst state from a list of states
//
// Helps combining an overall states, only based on a
// few numbers for various checks.
//
// Order of preference: Critical, Unknown, Warning, Ok
func WorstState(states ...int) int {
	overall := -1

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

	if overall < 0 || overall > 3 {
		overall = check.Unknown
	}

	return overall
}
