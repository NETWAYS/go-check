package perfdata

import (
	"fmt"
	"testing"

	"github.com/NETWAYS/go-check"
	"github.com/stretchr/testify/assert"
)

func ExamplePerfdata() {
	perf := Perfdata{
		Label: "test",
		Value: 10.1,
		Uom:   "%",
		Warn:  &check.Threshold{Upper: 80},
		Crit:  &check.Threshold{Upper: 90},
		Min:   0,
		Max:   100}

	fmt.Println(perf)

	// Output:
	// test=10.1%;80;90;0;100
}

func TestPerfdata(t *testing.T) {
	perf := Perfdata{
		Label: "test",
		Value: 2,
	}

	assert.Equal(t, "test=2", perf.String())
}

func TestPerfdata2(t *testing.T) {
	perf := Perfdata{
		Label: "test",
		Value: 2,
		Uom:   "%",
	}

	assert.Equal(t, "test=2%", perf.String())
}

func TestPerfdata3(t *testing.T) {
	perf := Perfdata{
		Label: "foo bar",
		Value: 2.76,
		Uom:   "m",
		Warn:  &check.Threshold{Lower: 10, Upper: 25, Inside: true},
		Crit:  &check.Threshold{Lower: 15, Upper: 20, Inside: false},
	}

	assert.Equal(t, "'foo bar'=2.76m;@10:25;15:20", perf.String())
}
