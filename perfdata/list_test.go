package perfdata

import "fmt"

func ExamplePerfdataList() {
	list := PerfdataList{}
	list.Add(&Perfdata{Label: "test1", Value: 23})
	list.Add(&Perfdata{Label: "test2", Value: 42})

	fmt.Println(list)

	// Output:
	// test1=23;;;; test2=42;;;;
}
