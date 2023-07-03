package convert

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBytes(t *testing.T) {
	b, err := ParseBytes("1024")
	assert.NoError(t, err)
	assert.IsType(t, BytesIEC(0), b)
	assert.Equal(t, uint64(1024), b.Bytes())

	b, err = ParseBytes("1MB")
	assert.NoError(t, err)
	assert.IsType(t, BytesSI(0), b)
	assert.Equal(t, uint64(1000*1000), b.Bytes())

	b, err = ParseBytes("1 MiB")
	assert.NoError(t, err)
	assert.IsType(t, BytesIEC(0), b)
	assert.Equal(t, uint64(1024*1024), b.Bytes())

	b, err = ParseBytes("100MB")
	assert.NoError(t, err)

	b, err = ParseBytes("100MiB")
	assert.NoError(t, err)

	b, err = ParseBytes("  23   GiB  ")
	assert.NoError(t, err)
	assert.IsType(t, BytesIEC(0), b)
	assert.Equal(t, uint64(23*1024*1024*1024), b.Bytes())

	b, err = ParseBytes("1.2.3.4MB")
	assert.Error(t, err)
	assert.Nil(t, b)

	b, err = ParseBytes("1PHD")
	assert.Error(t, err)
	assert.Nil(t, b)
}
