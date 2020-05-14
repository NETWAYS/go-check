package convert

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleParseBytes() {
	err, b := ParseBytes("100MB")
	if err != nil {
		panic(err)
	}

	fmt.Println(b)
	fmt.Println(b.ToHumanReadable())
	fmt.Println(b.ToKilobyte())
	fmt.Println(b.ToGigabyte())
	// Output: 100 MB
	// 100 MB
	// 100000
	// 0.1
}

func TestParseBytes(t *testing.T) {
	var b *Bytesize
	var err error

	err, b = ParseBytes(uint64(1))
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), b.Data)
	assert.Equal(t, "B", b.Unit)

	err, b = ParseBytes(1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), b.Data)
	assert.Equal(t, "B", b.Unit)

	err, b = ParseBytes("1")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), b.Data)
	assert.Equal(t, "B", b.Unit)

	err, b = ParseBytes("1mb")
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), b.Data)
	assert.Equal(t, "MB", b.Unit)

	err, b = ParseBytes("foobar")
	assert.Error(t, err)
}

func TestBytesize_String(t *testing.T) {
	err, b := ParseBytes("100mb")
	assert.NoError(t, err)
	assert.Equal(t, "100 MB", b.String())

	err, b = ParseBytes("100mb")
	assert.NoError(t, err)
	assert.Equal(t, "100 MB", fmt.Sprint(b))
}

func TestBytesize_ToHumanReadable(t *testing.T) {
	err, b := ParseBytes("900MB")
	assert.NoError(t, err)
	assert.Equal(t, "900 MB", b.ToHumanReadable())

	err, b = ParseBytes("1000mb")
	assert.NoError(t, err)
	assert.Equal(t, "1 GB", b.ToHumanReadable())
}
