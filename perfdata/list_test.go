package perfdata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePerfdataList() {
	list := PerfdataList{}
	list.Add(&Perfdata[int]{Label: "test1", Value: 23})
	list.Add(&Perfdata[float64]{Label: "test2", Value: 42.2})

	fmt.Println(list)

	// Output:
	// test1=23;;;0;0 test2=42.2;;;0;0
}

func TestPerfdataListFormating(t *testing.T) {
	list := PerfdataList{}
	list.Add(&Perfdata[int]{Label: "test1", Value: 23})
	list.Add(&Perfdata[float64]{Label: "test2", Value: 42.2})

	assert.Equal(t, "test1=23;;;0;0 test2=42.2;;;0;0", list.String())
}

func BenchmarkPerfdataListFormating(b *testing.B) {
	b.ReportAllocs()

	list := PerfdataList{}
	list.Add(&Perfdata[int]{Label: "test1", Value: 23})
	list.Add(&Perfdata[float64]{Label: "test2", Value: 42.2})

	for i := 0; i < b.N; i++ {
		l := list.String()
		_ = l
	}
}
