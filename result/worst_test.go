package result

import (
	"testing"

	"github.com/NETWAYS/go-check"
)

func TestWorstState2(t *testing.T) {
	if WorstState(check.Unknown) != check.Unknown {
		t.Fatalf("expected 3, got %d", WorstState(check.Unknown))
	}
	if WorstState(check.Critical) != check.Critical {
		t.Fatalf("expected 2, got %d", WorstState(check.Critical))
	}
	if WorstState(check.Warning) != check.Warning {
		t.Fatalf("expected 1, got %d", WorstState(check.Warning))
	}
	if WorstState(check.OK) != check.OK {
		t.Fatalf("expected 0, got %d", WorstState(check.OK))
	}
	if WorstState(check.OK, check.Warning, check.Critical, check.Unknown) != check.Critical {
		t.Fatalf("expected 2, got %d", WorstState(check.OK, check.Warning, check.Critical, check.Unknown))
	}
	if WorstState(check.OK, check.Warning, check.Unknown) != check.Unknown {
		t.Fatalf("expected 3, got %d", WorstState(check.OK, check.Warning, check.Unknown))
	}
	if WorstState(check.Warning, check.OK, check.OK) != check.Warning {
		t.Fatalf("expected 1, got %d", WorstState(check.Warning, check.OK, check.OK))
	}
	if WorstState(check.OK, check.OK, check.OK) != check.OK {
		t.Fatalf("expected 0, got %d", WorstState(check.OK, check.OK, check.OK))
	}

	// if WorstState(-1) != 3 {
	// 	t.Fatalf("expected 3, got %d", WorstState(-1))
	// }
	// if WorstState(4) != 3 {
	// 	t.Fatalf("expected 3, got %d", WorstState(4))
	// }
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
