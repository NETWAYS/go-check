package convert

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleParseBytes() {
	b := ParseBytes("100MB")

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

	b = ParseBytes(1)
	assert.Equal(t, 1, b.Data)
	assert.Equal(t, "B", b.Unit)

	b = ParseBytes("1")
	assert.Equal(t, 1, b.Data)
	assert.Equal(t, "B", b.Unit)

	b = ParseBytes("1mb")
	assert.Equal(t, 1, b.Data)
	assert.Equal(t, "MB", b.Unit)

	t.Skip("missing error handling")
	b = ParseBytes("foobar")
	assert.Equal(t, 1, b.Data)
	assert.Equal(t, "MB", b.Unit)
}

func TestBytesize_String(t *testing.T) {
	assert.Equal(t, "100 MB", ParseBytes("100mb").String())
	assert.Equal(t, "100 MB", fmt.Sprint(ParseBytes("100mb")))
}

func TestBytesize_ToHumanReadable(t *testing.T) {
	assert.Equal(t, "900 MB", ParseBytes("900MB").ToHumanReadable())
	assert.Equal(t, "1 GB", ParseBytes("1000MB").ToHumanReadable())
	assert.Equal(t, "1.2 GB", ParseBytes("1200MB").ToHumanReadable())
}
