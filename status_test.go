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

func TestStatusInt(t *testing.T) {
	testcases := map[string]struct {
		input    string
		expected int
	}{
		"OK": {
			expected: 0,
			input:    "OK",
		},
		"WARNING": {
			expected: 1,
			input:    "warning",
		},
		"CRITICAL": {
			expected: 2,
			input:    "Critical",
		},
		"UNKNOWN": {
			expected: 3,
			input:    "unknown",
		},
		"Invalid-Input": {
			expected: 3,
			input:    "Something else",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := StatusInt(tc.input)

			if actual != tc.expected {
				t.Error("\nActual: ", actual, "\nExpected: ", tc.expected)
			}
		})
	}
}
