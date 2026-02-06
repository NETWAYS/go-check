package check

import (
	"testing"
)

func TestStatus_String(t *testing.T) {
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

func TestStatus_FromString(t *testing.T) {
	testcases := map[string]struct {
		expected Status
		input    string
	}{
		"OK": {
			input:    "OK",
			expected: OK,
		},
		"WARNING": {
			input:    "WARNING",
			expected: Warning,
		},
		"CRITICAL": {
			input:    "CRITICAL",
			expected: Critical,
		},
		"UNKNOWN": {
			input:    "UNKNOWN",
			expected: Unknown,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual, _ := NewStatusFromString(tc.input)

			if actual != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestStatus_FromInt(t *testing.T) {
	testcases := map[string]struct {
		expected Status
		input    int
	}{
		"OK": {
			input:    0,
			expected: OK,
		},
		"WARNING": {
			input:    1,
			expected: Warning,
		},
		"CRITICAL": {
			input:    2,
			expected: Critical,
		},
		"UNKNOWN": {
			input:    3,
			expected: Unknown,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual, _ := NewStatus(tc.input)

			if actual != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
