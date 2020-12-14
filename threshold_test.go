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

func TestBoundToString(t *testing.T) {
	assert.Equal(t, "10", BoundToString(10))
	assert.Equal(t, "10.1", BoundToString(10.1))
	assert.Equal(t, "10.0000000000001", BoundToString(10.0000000000001))
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
