package check

import (
	"testing"
)

func TestStatusText(t *testing.T) {
	testcases := map[string]struct {
		input    int
		expected string
	}{
		"OK": {
			input:    0,
			expected: "OK",
		},
		"WARNING": {
			input:    1,
			expected: "WARNING",
		},
		"CRITICAL": {
			input:    2,
			expected: "CRITICAL",
		},
		"UNKNOWN": {
			input:    3,
			expected: "UNKNOWN",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := StatusText(tc.input)

			if actual != tc.expected {
				t.Error("\nActual: ", actual, "\nExpected: ", tc.expected)
			}
		})
	}
}
