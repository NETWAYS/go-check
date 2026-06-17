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
		"WaRnInG": {
			input:    "WaRnInG",
			expected: Warning,
		},
		"CRITICAL": {
			input:    "CRITICAL",
			expected: Critical,
		},
		"Critical": {
			input:    "Critical",
			expected: Critical,
		},
		"UNKNOWN": {
			input:    "UNKNOWN",
			expected: Unknown,
		},
		"unknown": {
			input:    "unknown",
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

func TestStatus_FromString_WithErr(t *testing.T) {
	actual, err := NewStatusFromString("unittest")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if actual != Unknown {
		t.Fatalf("expected Unknown, got %v", actual)
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

func TestStatus_FromInt_WithErr(t *testing.T) {
	actual, err := NewStatus(1337)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if actual != Unknown {
		t.Fatalf("expected Unknown, got %v", actual)
	}
}

type CompareSet struct {
	Left   Status
	Right  Status
	Result int
}

func TestCompareStatus(t *testing.T) {
	inputList := []Status{OK, Warning, Unknown, Critical}

	for _, val := range inputList {
		if Compare(val, val) != 0 {
			t.Fatalf("Equal comparison failed with %d and %d", val, val)
		}
	}

	ComparisonTests := []CompareSet{
		{OK, Critical, 1},
		{OK, Warning, 1},
		{OK, Unknown, 1},
		{Warning, Critical, 1},
		{Warning, Unknown, 1},
		{Warning, OK, -1},
		{Unknown, Critical, 1},
		{Unknown, OK, -1},
		{Unknown, Warning, -1},
		{Critical, OK, -1},
		{Critical, Warning, -1},
		{Critical, Unknown, -1},
	}

	for _, set := range ComparisonTests {
		if Compare(set.Left, set.Right) != set.Result {
			t.Fatalf("Comparison failed with %s and %s, got %d", set.Left, set.Right, Compare(set.Left, set.Right))
		}
	}
}

func TestWorstState(t *testing.T) {
	tests := []struct {
		name     string
		input    []Status
		expected Status
	}{
		{
			name:     "Unknown",
			input:    []Status{Unknown},
			expected: Unknown,
		},
		{
			name:     "Unknown",
			input:    []Status{Unknown},
			expected: Unknown,
		},
		{
			name:     "Critical",
			input:    []Status{Critical},
			expected: Critical,
		},
		{
			name:     "Warning",
			input:    []Status{Warning},
			expected: Warning,
		},
		{
			name:     "OK",
			input:    []Status{OK},
			expected: OK,
		},
		{
			name:     "Mixed with Critical",
			input:    []Status{OK, Warning, Critical, Unknown},
			expected: Critical,
		},
		{
			name:     "Mixed order with Critical",
			input:    []Status{OK, Critical, Warning, Unknown},
			expected: Critical,
		},
		{
			name:     "Mixed with Unknown",
			input:    []Status{OK, Warning, Unknown},
			expected: Unknown,
		},
		{
			name:     "Mixed order with Unknown",
			input:    []Status{OK, Unknown, Warning},
			expected: Unknown,
		},
		{
			name:     "Warning with OKs",
			input:    []Status{Warning, OK, OK},
			expected: Warning,
		},
		{
			name:     "Warning orderwith OKs",
			input:    []Status{OK, Warning, OK},
			expected: Warning,
		},
		{
			name:     "All OK",
			input:    []Status{OK, OK, OK},
			expected: OK,
		},
		{
			name:     "Empty",
			input:    []Status{},
			expected: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WorstState(tt.input...); got != tt.expected {
				t.Errorf("WorstState(%v) is: %v, expected %v", tt.input, got, tt.expected)
			}
		})
	}
}

func BenchmarkWorstState(b *testing.B) {
	b.ReportAllocs()

	// Initialize slice for benchmarking
	states := make([]Status, 0, 100)
	for i := range 100 {
		st, _ := NewStatus(i % 4)
		states = append(states, st)
	}

	for i := 0; i < b.N; i++ {
		s := WorstState(states...)
		_ = s
	}
}
