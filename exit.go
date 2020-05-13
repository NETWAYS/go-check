package check

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/mitchellh/go-ps"
)

// AllowExit lets you disable the call to os.Exit() in ExitXxx() functions of this package.
//
// This should be used carefully and most likely only for testing.
var AllowExit = true

// PrintStack prints the error stack when recovering from a panic with CatchPanic()
var PrintStack = true

// Exit prints the plugin output and exits the program
func Exit(rc int, output string, args ...interface{}) {
	fmt.Println(StatusText(rc), "-", fmt.Sprintf(output, args...))
	BaseExit(rc)
}

func BaseExit(rc int) {
	if AllowExit {
		os.Exit(rc)
	} else {
		fmt.Println("would exit with code", rc)
	}
}

// ExitError exists with an Unknown state while reporting the error
//
// TODO: more information about the error
func ExitError(err error) {
	Exit(Unknown, err.Error())
}

// CatchPanic is a general function for defer, to capture any panic that occurred during runtime of a check
//
// The function will recover from the condition and exit with a proper UNKNOWN status, while showing error
// and the call stack.
func CatchPanic() {
	ppid := os.Getppid()
	if parent, err := ps.FindProcess(ppid); err == nil {
		if parent.Executable() == "dlv" {
			// seems to be a debugger, don't continue with recover
			return
		}
	}

	if r := recover(); r != nil {
		output := fmt.Sprint("Golang encountered a panic: ", r)
		if PrintStack {
			output += "\n\n" + string(debug.Stack())
		}

		Exit(Unknown, output)
	}
}
