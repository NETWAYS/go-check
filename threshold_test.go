package check

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testThresholds = map[string]*Threshold{
	"10":     {Lower: 0, Upper: 10},
	"10:":    {Lower: 10, Upper: PosInf},
	"~:10":   {Lower: NegInf, Upper: 10},
	"10:20":  {Lower: 10, Upper: 20},
	"@10:20": {Lower: 10, Upper: 20, Inside: true},
	"-10:10": {Lower: -10, Upper: 10},
	"":       nil,
}

func TestBoundaryToString(t *testing.T) {
	assert.Equal(t, "10", BoundaryToString(10))
	assert.Equal(t, "10.1", BoundaryToString(10.1))
	assert.Equal(t, "10.0000000000001", BoundaryToString(10.0000000000001))
}

func TestParseThreshold(t *testing.T) {
	for spec, ref := range testThresholds {
		th, err := ParseThreshold(spec)

		if ref == nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, ref, th)
		}
	}
}

func TestThreshold_String(t *testing.T) {
	for spec, ref := range testThresholds {
		if ref != nil {
			assert.Equal(t, spec, ref.String())
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
	assert.NoError(t, err)
	assert.True(t, thr.DoesViolate(11))
	assert.False(t, thr.DoesViolate(10))
	assert.False(t, thr.DoesViolate(0))
	assert.True(t, thr.DoesViolate(-1))

	thr, err = ParseThreshold("10:")
	assert.NoError(t, err)
	assert.False(t, thr.DoesViolate(3000))
	assert.False(t, thr.DoesViolate(10))
	assert.True(t, thr.DoesViolate(9))
	assert.True(t, thr.DoesViolate(0))
	assert.True(t, thr.DoesViolate(-1))

	thr, err = ParseThreshold("~:10")
	assert.NoError(t, err)
	assert.False(t, thr.DoesViolate(-3000))
	assert.False(t, thr.DoesViolate(0))
	assert.False(t, thr.DoesViolate(10))
	assert.True(t, thr.DoesViolate(11))
	assert.True(t, thr.DoesViolate(3000))

	thr, err = ParseThreshold("10:20")
	assert.NoError(t, err)
	assert.False(t, thr.DoesViolate(10))
	assert.False(t, thr.DoesViolate(15))
	assert.False(t, thr.DoesViolate(20))
	assert.True(t, thr.DoesViolate(9))
	assert.True(t, thr.DoesViolate(-1))
	assert.True(t, thr.DoesViolate(20.1))
	assert.True(t, thr.DoesViolate(3000))

	thr, err = ParseThreshold("@10:20")
	assert.NoError(t, err)
	assert.True(t, thr.DoesViolate(10))
	assert.True(t, thr.DoesViolate(15))
	assert.True(t, thr.DoesViolate(20))
	assert.False(t, thr.DoesViolate(9))
	assert.False(t, thr.DoesViolate(-1))
	assert.False(t, thr.DoesViolate(20.1))
	assert.False(t, thr.DoesViolate(3000))
}
