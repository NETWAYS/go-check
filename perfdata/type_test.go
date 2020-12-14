package perfdata

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExamplePerfdata() {
	perf := Perfdata{
		Label: "test",
		Value: 10.1,
		Uom:   "%",
		Warn:  &Threshold{Upper: 80},
		Crit:  &Threshold{Upper: 90},
		Min:   0, Max: 100}

	fmt.Println(perf)

	// Output:
	// test=10.1%;80;90;0;100
}

func TestFormatNumeric(t *testing.T) {
	assert.Equal(t, "10", FormatNumeric(10))
	assert.Equal(t, "-10", FormatNumeric(-10))
	assert.Equal(t, "10", FormatNumeric(uint8(10)))
	assert.Equal(t, "1234.5678", FormatNumeric(1234.5678))
	assert.Equal(t, "1234.567", FormatNumeric(float32(1234.567)))
	assert.Equal(t, "1234.567", FormatNumeric("1234.567"))
}

func TestFormatLabel(t *testing.T) {
	assert.Equal(t, "test", FormatLabel("test"))
	assert.Equal(t, "'test test'", FormatLabel("test test"))
	assert.Equal(t, "test_x", FormatLabel("test\t\n\\x"))
}
