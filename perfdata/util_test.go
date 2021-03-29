package perfdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	assert.Equal(t, "t_est_x", FormatLabel("t%$%^est\t\n\\x"))
	assert.Equal(t, "test/:x", FormatLabel("test/:x"))
}

func TestIsValidUom(t *testing.T) {
	assert.True(t, IsValidUom("%"))
	assert.False(t, IsValidUom(""))
	assert.False(t, IsValidUom("X"))
}
