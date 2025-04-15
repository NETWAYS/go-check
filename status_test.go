package check

import (
	"testing"
)

func TestStatusText(t *testing.T) {
	testcases := map[string]struct {
		input    Status
		expected string
	}{
		"OK": {
			input:    OK,
			expected: "OK",
		},
		"WARNING": {
			input:    Warning,
			expected: "WARNING",
		},
		"CRITICAL": {
			input:    Critical,
			expected: "CRITICAL",
		},
		"UNKNOWN": {
			input:    Unknown,
			expected: "UNKNOWN",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {

			actual := tc.input.String()

			if actual != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
