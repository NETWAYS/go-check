package perfdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatNumeric(t *testing.T) {
	assert.Equal(t, "10", FormatNumeric(10))
	assert.Equal(t, "-10", FormatNumeric(-10))
	assert.Equal(t, "10", FormatNumeric(uint8(10)))
	assert.Equal(t, "1234.567", FormatNumeric(1234.567))
	assert.Equal(t, "1234.567", FormatNumeric(float32(1234.567)))
	assert.Equal(t, "1234.567", FormatNumeric("1234.567"))
	assert.Equal(t, "1234567890.988", FormatNumeric(1234567890.9877))
}

func TestFormatLabel(t *testing.T) {
	assert.Equal(t, "test", FormatLabel("test"))
	assert.Equal(t, "'test test'", FormatLabel("test test"))
	assert.Equal(t, "test_x", FormatLabel("test\t\n\\x"))
	assert.Equal(t, "t_est_x", FormatLabel("t%$%^est\t\n\\x"))
	assert.Equal(t, "test/:x", FormatLabel("test/:x"))
}
