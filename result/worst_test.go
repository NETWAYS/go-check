package result

import (
	"testing"

	"github.com/NETWAYS/go-check"
)

func TestWorstState2(t *testing.T) {
	if WorstState(3) != 3 {
		t.Fatalf("expected 3, got %d", WorstState(3))
	}
	if WorstState(2) != 2 {
		t.Fatalf("expected 2, got %d", WorstState(2))
	}
	if WorstState(1) != 1 {
		t.Fatalf("expected 1, got %d", WorstState(1))
	}
	if WorstState(0) != 0 {
		t.Fatalf("expected 0, got %d", WorstState(0))
	}

	if WorstState(0, 1, 2, 3) != 2 {
		t.Fatalf("expected 2, got %d", WorstState(0, 1, 2, 3))
	}
	if WorstState(0, 1, 3) != 3 {
		t.Fatalf("expected 3, got %d", WorstState(0, 1, 3))
	}
	if WorstState(1, 0, 0) != 1 {
		t.Fatalf("expected 1, got %d", WorstState(1, 0, 0))
	}
	if WorstState(0, 0, 0) != 0 {
		t.Fatalf("expected 0, got %d", WorstState(0, 0, 0))
	}

	if WorstState(-1) != 3 {
		t.Fatalf("expected 3, got %d", WorstState(-1))
	}
	if WorstState(4) != 3 {
		t.Fatalf("expected 3, got %d", WorstState(4))
	}
}

func BenchmarkWorstState(b *testing.B) {
	b.ReportAllocs()

	// Initialize slice for benchmarking
	states := make([]check.Status, 0, 100)
	for i := 0; i < 100; i++ {
		s, _ := check.NewStatusFromInt(i % 4)
		states = append(states, s)
	}

	for i := 0; i < b.N; i++ {
		s := WorstState(states...)
		_ = s
	}
}
