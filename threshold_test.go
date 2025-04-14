package check

import (
	"math"
	"reflect"
	"testing"
)

var testThresholds = map[string]*Threshold{
	"10":             {Lower: 0, Upper: 10},
	"10:":            {Lower: 10, Upper: PosInf},
	"~:10":           {Lower: NegInf, Upper: 10},
	"10:20":          {Lower: 10, Upper: 20},
	"@10:20":         {Lower: 10, Upper: 20, Inside: true},
	"-10:10":         {Lower: -10, Upper: 10},
	"-10.001:10.001": {Lower: -10.001, Upper: 10.001},
	"":               nil,
}

func TestBoundaryToString(t *testing.T) {
	if BoundaryToString(10) != "10" {
		t.Fatalf("expected '10', got %s", BoundaryToString(10))
	}
	if BoundaryToString(10.1) != "10.1" {
		t.Fatalf("expected '10.1', got %s", BoundaryToString(10.1))
	}
	if BoundaryToString(10.001) != "10.001" {
		t.Fatalf("expected '10.001', got %s", BoundaryToString(10.001))
	}
}

func TestParseThreshold(t *testing.T) {
	for spec, ref := range testThresholds {
		th, err := ParseThreshold(spec)

		if ref == nil {
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		} else {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(ref, th) {
				t.Fatalf("expected %v, got %v for spec %s", ref, th, spec)
			}
			if th.String() != spec {
				t.Fatalf("expected %s, got %s for spec %s", spec, th.String(), spec)
			}
		}
	}
}

func TestThreshold_String(t *testing.T) {
	for spec, ref := range testThresholds {
		if ref != nil {
			if spec != ref.String() {
				t.Fatalf("expected %v, got %v", ref.String(), spec)
			}
		}
	}
}

// Threshold  Generate an alert if x...
// 10         < 0 or > 10, (outside the range of {0 .. 10})
// 10:        < 10, (outside {10 .. ∞})
// ~:10       > 10, (outside the range of {-∞ .. 10})
// 10:20      < 10 or > 20, (outside the range of {10 .. 20})
// @10:20     ≥ 10 and ≤ 20, (inside the range of {10 .. 20})
func TestThreshold_DoesViolate(t *testing.T) {
	thr, err := ParseThreshold("10")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !thr.DoesViolate(11) {
		t.Fatalf("expected true, got false")
	}
	if thr.DoesViolate(10) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(0) {
		t.Fatalf("expected false, got true")
	}
	if !thr.DoesViolate(-1) {
		t.Fatalf("expected true, got false")
	}

	thr, err = ParseThreshold("10:")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if thr.DoesViolate(3000) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(10) {
		t.Fatalf("expected false, got true")
	}
	if !thr.DoesViolate(9) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(0) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(-1) {
		t.Fatalf("expected true, got false")
	}

	thr, err = ParseThreshold("~:10")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if thr.DoesViolate(-3000) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(0) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(10) {
		t.Fatalf("expected false, got true")
	}
	if !thr.DoesViolate(11) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(3000) {
		t.Fatalf("expected true, got false")
	}

	thr, err = ParseThreshold("10:20")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if thr.DoesViolate(10) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(15) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(20) {
		t.Fatalf("expected false, got true")
	}
	if !thr.DoesViolate(9) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(-1) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(20.1) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(3000) {
		t.Fatalf("expected true, got false")
	}

	thr, err = ParseThreshold("@10:20")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !thr.DoesViolate(10) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(15) {
		t.Fatalf("expected true, got false")
	}
	if !thr.DoesViolate(20) {
		t.Fatalf("expected true, got false")
	}
	if thr.DoesViolate(9) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(-1) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(20.1) {
		t.Fatalf("expected false, got true")
	}
	if thr.DoesViolate(3000) {
		t.Fatalf("expected false, got true")
	}
}

func TestFormatFloat(t *testing.T) {
	if FormatFloat(1000000000000) != "1000000000000" {
		t.Fatalf("expected '1000000000000', got %s", FormatFloat(1000000000000))
	}
	if FormatFloat(1000000000) != "1000000000" {
		t.Fatalf("expected '1000000000', got %s", FormatFloat(1000000000))
	}
	if FormatFloat(1234567890.9877) != "1234567890.988" {
		t.Fatalf("expected '1234567890.988', got %s", FormatFloat(1234567890.9877))
	}
	if FormatFloat(math.Inf(-1)) != "-Inf" {
		t.Fatalf("expected '-Inf', got %s", FormatFloat(math.Inf(-1)))
	}
	if FormatFloat(math.Inf(1)) != "+Inf" {
		t.Fatalf("expected '+Inf', got %s", FormatFloat(math.Inf(1)))
	}
}
