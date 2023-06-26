package perfdata

import (
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/stretchr/testify/assert"
)

func BenchmarkPerfdataString(b *testing.B) {
	b.ReportAllocs()

	perf := Perfdata{
		Label: "test test=test",
		Value: 10.1,
		Uom:   "%",
		Warn:  &check.Threshold{Upper: 80},
		Crit:  &check.Threshold{Upper: 90},
		Min:   0,
		Max:   100}

	for i := 0; i < b.N; i++ {
		p := perf.String()
		_ = p
	}
}

func TestRenderPerfdata(t *testing.T) {
	testcases := map[string]struct {
		perf     Perfdata
		expected string
	}{
		"simple": {
			perf: Perfdata{
				Label: "test",
				Value: 2,
			},
			expected: "test=2",
		},
		"with-special-chars": {
			perf: Perfdata{
				Label: "test_🖥️_'test",
				Value: 2,
			},
			expected: "test_🖥️__test=2",
		},
		"with-uom": {
			perf: Perfdata{
				Label: "test",
				Value: 2,
				Uom:   "%",
			},
			expected: "test=2%",
		},
		"with-thresholds": {
			perf: Perfdata{
				Label: "foo bar",
				Value: 2.76,
				Uom:   "m",
				Warn:  &check.Threshold{Lower: 10, Upper: 25, Inside: true},
				Crit:  &check.Threshold{Lower: 15, Upper: 20, Inside: false},
			},
			expected: "'foo bar'=2.76m;@10:25;15:20",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.perf.String())
		})
	}
}

func TestFormatNumeric(t *testing.T) {
	assert.Equal(t, "10", formatNumeric(10))
	assert.Equal(t, "-10", formatNumeric(-10))
	assert.Equal(t, "10", formatNumeric(uint8(10)))
	assert.Equal(t, "1234.567", formatNumeric(1234.567))
	assert.Equal(t, "1234.567", formatNumeric(float32(1234.567)))
	assert.Equal(t, "1234.567", formatNumeric("1234.567"))
	assert.Equal(t, "1234567890.988", formatNumeric(1234567890.9877))
}
