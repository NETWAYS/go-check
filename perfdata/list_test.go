package perfdata

import (
	"fmt"
	"testing"
)

func ExamplePerfdataList() {
	list := PerfdataList{}
	list.Add(&Perfdata{Label: "test1", Value: 23})
	list.Add(&Perfdata{Label: "test2", Value: 42})

	fmt.Println(list)

	// Output:
	// [test1=23 test2=42]
}

func TestPerfdataListFormating(t *testing.T) {
	list := PerfdataList{}
	list.Add(&Perfdata{Label: "test1", Value: 23})
	list.Add(&Perfdata{Label: "test2", Value: 42})

	actual := list.String()
	expected := "test1=23 test2=42"

	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func BenchmarkPerfdataListFormating(b *testing.B) {
	b.ReportAllocs()

	list := PerfdataList{}
	list.Add(&Perfdata{Label: "test1", Value: 23})
	list.Add(&Perfdata{Label: "test2", Value: 42})

	for i := 0; i < b.N; i++ {
		l := list.String()
		_ = l
	}
}
