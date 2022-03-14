package convert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesSI_HumanReadable(t *testing.T) {
	assert.Equal(t, "0B", BytesSI(0).HumanReadable())
	assert.Equal(t, "999B", BytesSI(999).HumanReadable())
	assert.Equal(t, "999KB", BytesSI(999*1000).HumanReadable())
	assert.Equal(t, "999MB", BytesSI(999*1000*1000).HumanReadable())
	assert.Equal(t, "999GB", BytesSI(999*1000*1000*1000).HumanReadable())
	assert.Equal(t, "999TB", BytesSI(999*1000*1000*1000*1000).HumanReadable())

	assert.Equal(t, "4PB", BytesSI(4*1000*1000*1000*1000*1000).HumanReadable())
	assert.Equal(t, "4000PB", BytesSI(4*1000*1000*1000*1000*1000*1000).HumanReadable())

	assert.Equal(t, "4TB", BytesSI(4*1000*1000*1000*1000).HumanReadable())
	assert.Equal(t, "4PB", BytesSI(4*1000*1000*1000*1000*1000).HumanReadable())

	assert.Equal(t, "1263MB", BytesSI(1263*1000*1000).HumanReadable()) // and not 1.23GiB

	assert.Equal(t, "123.05MB", BytesSI(123050*1000).HumanReadable())
	assert.Equal(t, "14.67GB", BytesSI(14670*1000*1000).HumanReadable())

	assert.Equal(t, "1000KB", BytesSI(1000*1000).HumanReadable())
	assert.Equal(t, "2MB", BytesSI(2*1000*1000).HumanReadable())
	assert.Equal(t, "3MB", BytesSI(3*1000*1000).HumanReadable())
}
