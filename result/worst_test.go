package result

import (
	"testing"

	"github.com/NETWAYS/go-check"
)

func TestWorstState(t *testing.T) {
	tests := []struct {
		name     string
		input    []check.Status
		expected check.Status
	}{
		{
			name:     "Unknown",
			input:    []check.Status{check.Unknown},
			expected: check.Unknown,
		},
		{
			name:     "Unknown",
			input:    []check.Status{check.Unknown},
			expected: check.Unknown,
		},
		{
			name:     "Critical",
			input:    []check.Status{check.Critical},
			expected: check.Critical,
		},
		{
			name:     "Warning",
			input:    []check.Status{check.Warning},
			expected: check.Warning,
		},
		{
			name:     "OK",
			input:    []check.Status{check.OK},
			expected: check.OK,
		},
		{
			name:     "Mixed with Critical",
			input:    []check.Status{check.OK, check.Warning, check.Critical, check.Unknown},
			expected: check.Critical,
		},
		{
			name:     "Mixed order with Critical",
			input:    []check.Status{check.OK, check.Critical, check.Warning, check.Unknown},
			expected: check.Critical,
		},
		{
			name:     "Mixed with Unknown",
			input:    []check.Status{check.OK, check.Warning, check.Unknown},
			expected: check.Unknown,
		},
		{
			name:     "Mixed order with Unknown",
			input:    []check.Status{check.OK, check.Unknown, check.Warning},
			expected: check.Unknown,
		},
		{
			name:     "Warning with OKs",
			input:    []check.Status{check.Warning, check.OK, check.OK},
			expected: check.Warning,
		},
		{
			name:     "Warning orderwith OKs",
			input:    []check.Status{check.OK, check.Warning, check.OK},
			expected: check.Warning,
		},
		{
			name:     "All OK",
			input:    []check.Status{check.OK, check.OK, check.OK},
			expected: check.OK,
		},
		{
			name:     "Empty",
			input:    []check.Status{},
			expected: check.Unknown,
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
	states := make([]check.Status, 0, 100)
	for i := range 100 {
		st, _ := check.NewStatus(i % 4)
		states = append(states, st)
	}

	for i := 0; i < b.N; i++ {
		s := WorstState(states...)
		_ = s
	}
}
