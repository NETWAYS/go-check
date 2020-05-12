package check

import (
	"fmt"
)

func ExampleStatusText() {
	fmt.Println(StatusText(OK), StatusText(Warning), StatusText(Critical), StatusText(Unknown))
	// Output: OK WARNING CRITICAL UNKNOWN
}
