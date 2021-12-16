package convert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesIEC_HumanReadable(t *testing.T) {
	assert.Equal(t, "999B", BytesIEC(999).HumanReadable())
	assert.Equal(t, "999KiB", BytesIEC(999*1024).HumanReadable())
	assert.Equal(t, "999MiB", BytesIEC(999*1024*1024).HumanReadable())
	assert.Equal(t, "999GiB", BytesIEC(999*1024*1024*1024).HumanReadable())
	assert.Equal(t, "999TiB", BytesIEC(999*1024*1024*1024*1024).HumanReadable())
	assert.Equal(t, "4PiB", BytesIEC(4*1024*1024*1024*1024*1024).HumanReadable())
	assert.Equal(t, "4096PiB", BytesIEC(4*1024*1024*1024*1024*1024*1024).HumanReadable())

	assert.Equal(t, "1263MiB", BytesIEC(1263*1024*1024).HumanReadable()) // and not 1.23GiB

	assert.Equal(t, "100MiB", BytesIEC(100*1024*1024).HumanReadable())

	assert.Equal(t, "123.05MiB", BytesIEC(129032519).HumanReadable())
	assert.Equal(t, "14.67GiB", BytesIEC(15756365824).HumanReadable())

	assert.Equal(t, "1024KiB", BytesIEC(1024*1024).HumanReadable())
	assert.Equal(t, "2MiB", BytesIEC(2*1024*1024).HumanReadable())
}
