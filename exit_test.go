package check

import (
	"fmt"
	"os"
	"testing"
)

func ExampleExit() {
	Exit(OK, "Everything is fine")
	// Output: OK - Everything is fine
	// would exit with code 0
}

func ExampleExitError() {
	err := fmt.Errorf("connection to %s has been timed out", "localhost:12345")
	ExitError(err)
	// Output: UNKNOWN - connection to localhost:12345 has been timed out
	// would exit with code 3
}

func ExampleCatchPanic() {
	defer CatchPanic()

	panic("something bad happened")
	// Output: UNKNOWN - Golang encountered a panic: something bad happened
	// would exit with code 3
}

func TestMain(m *testing.M) {
	// disable actual exit
	AllowExit = false

	// disable stack trace for the example
	PrintStack = false

	os.Exit(m.Run())
}
