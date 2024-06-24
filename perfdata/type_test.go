package perfdata

import (
	"math"
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/stretchr/testify/assert"
)

func BenchmarkPerfdataString(b *testing.B) {
	b.ReportAllocs()

	perf := Perfdata[float64]{
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
		perf     Perfdata[int]
		expected string
	}{
		"simple": {
			perf: Perfdata[int]{
				Label: "test",
				Value: 2,
			},
			expected: "test=2;;;0;0",
		},
		"with-quotes": {
			perf: Perfdata[int]{
				Label: "te's\"t",
				Value: 2,
			},
			expected: "te_s_t=2;;;0;0",
		},
		"with-special-chars": {
			perf: Perfdata[int]{
				Label: "test_🖥️_'test",
				Value: 2,
			},
			expected: "test_🖥️__test=2;;;0;0",
		},
		"with-uom": {
			perf: Perfdata[int]{
				Label: "test",
				Value: 2,
				Uom:   "%",
			},
			expected: "test=2%;;;0;0",
		},
		"with-thresholds": {
			perf: Perfdata[int]{
				Label: "foo bar",
				Value: 2,
				Uom:   "m",
				Warn:  &check.Threshold{Lower: 10, Upper: 25, Inside: true},
				Crit:  &check.Threshold{Lower: 15, Upper: 20, Inside: false},
			},
			expected: "'foo bar'=2m;@10:25;15:20;0;0",
		},
	}

	testcasesWithErrors := map[string]struct {
		perf     Perfdata[float64]
		expected string
	}{
		"invalid-value": {
			perf: Perfdata[float64]{
				Label: "to infinity",
				Value: math.Inf(+1),
			},
			expected: "",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			pfVal, err := tc.perf.ValidatedString()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, pfVal)
		})
	}

	for name, tc := range testcasesWithErrors {
		t.Run(name, func(t *testing.T) {
			pfVal, err := tc.perf.ValidatedString()
			assert.Error(t, err)
			assert.Equal(t, tc.expected, pfVal)
		})
	}
}

type pfFormatTest struct {
	Result     string
	InputValue interface{}
}

func TestFormatNumeric(t *testing.T) {
	testdata := []pfFormatTest{
		{
			Result:     "10",
			InputValue: 10,
		},
		{
			Result:     "-10",
			InputValue: -10,
		},
		{
			Result:     "10",
			InputValue: uint8(10),
		},
		{
			Result:     "1234.567",
			InputValue: float64(1234.567),
		},
		{
			Result:     "3456.789",
			InputValue: float32(3456.789),
		},
		{
			Result:     "1234567890.988",
			InputValue: 1234567890.9877,
		},
	}

	for _, val := range testdata {
		formatted, err := formatNumeric(val.InputValue)
		assert.NoError(t, err)
		assert.Equal(t, val.Result, formatted)
	}
}
