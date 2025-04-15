package result

import "github.com/NETWAYS/go-check"

// Determines the worst state from a list of states
//
// Helps combining an overall states, only based on a
// few numbers for various checks.
//
// Order of preference: Critical, Unknown, Warning, Ok
func WorstState(states ...check.Status) check.Status {
	overall := -1

	for _, state := range states {
		if state == check.Critical {
			overall = int(check.Critical)
		} else if state == check.Unknown {
			if overall != int(check.Critical) {
				overall = int(check.Unknown)
			}
		} else if int(state) > overall {
			overall = int(state)
		}
	}

	if overall < 0 || overall > 3 {
		overall = int(check.Unknown)
	}

	s, _ := check.NewStatusFromInt(overall)

	return s
}
