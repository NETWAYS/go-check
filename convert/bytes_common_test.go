package convert

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleParseBytes() {
	bSI, err := ParseBytes("100MB")
	if err != nil {
		log.Fatal(err)
	}

	bIEC, err := ParseBytes("100MiB")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SI units")
	fmt.Println(bSI.Bytes())
	fmt.Println(bSI)
	fmt.Println(bSI.HumanReadable())

	fmt.Println("IEC units")
	fmt.Println(bIEC.Bytes())
	fmt.Println(bIEC)
	fmt.Println(bIEC.HumanReadable())
	// Output:
	// SI units
	// 100000000
	// 100MB
	// 100MB
	// IEC units
	// 104857600
	// 100MiB
	// 100MiB
}

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
