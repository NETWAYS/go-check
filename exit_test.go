package check

import (
	"fmt"
	"os"
	"testing"
)

func ExampleExit() {
	Exit(OK, fmt.Sprintf("Everything is fine - value=%d", 42))
	// Output: [OK] - Everything is fine - value=42
	// would exit with code 0
}

func ExampleExit_withVerticalBar() {
	Exit(Warning, fmt.Sprintf("Something|is|wrong - value=%d", 42))
	// Output: [WARNING] - Something is wrong - value=42
	// would exit with code 1
}

// TODO: How to handle this?
func ExampleExit_withPerfdata() {
	Exit(Critical, "Everything is broken", "|", "percent_packet_loss=100")
	// Output: [CRITICAL] - Everything is broken | percent_packet_loss=100
	// would exit with code 2
}

func ExampleExitError() {
	err := fmt.Errorf("connection to %s has been timed out", "localhost:12345")
	ExitError(err)
	// Output: [UNKNOWN] - connection to localhost:12345 has been timed out (*errors.errorString)
	// would exit with code 3
}

func ExampleCatchPanic() {
	defer CatchPanic()

	panic("something bad happened")
	// Output: [UNKNOWN] - Golang encountered a panic: something bad happened
	// would exit with code 3
}

func TestMain(m *testing.M) {
	// disable actual exit
	AllowExit = false

	// disable stack trace for the example
	PrintStack = false

	os.Exit(m.Run())
}
