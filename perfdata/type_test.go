package perfdata

import (
	"fmt"
	"github.com/NETWAYS/go-check"
)

func ExamplePerfdata() {
	perf := Perfdata{
		Label: "test",
		Value: 10.1,
		Uom:   "%",
		Warn:  &check.Threshold{Upper: 80},
		Crit:  &check.Threshold{Upper: 90},
		Min:   0, Max: 100}

	fmt.Println(perf)

	// Output:
	// test=10.1%;80;90;0;100
}
