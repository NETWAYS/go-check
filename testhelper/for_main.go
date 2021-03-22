package testhelper

import (
	"bytes"
	"github.com/NETWAYS/go-check"
	"io"
	"os"
)

// Execute main function from a main package, while capturing its stdout
//
// You will need to pass `main`, since the function won't have access to the package's namespace
func RunMainTest(f func(), args ...string) string {
	base := []string{"check_with_go_test"}
	origArgs := os.Args

	os.Args = append(base, args...)
	stdout := CaptureStdout(f)
	os.Args = origArgs

	return stdout
}

// Enable test mode by disabling the default exit behavior, so go test won't fail with plugin states
func EnableTestMode() {
	// disable actual exit
	check.AllowExit = false

	// disable stack trace for the example
	check.PrintStack = false
}

// Disable test mode behavior again
//
// Optional after testing has been done
func DisableTestMode() {
	// disable actual exit
	check.AllowExit = true

	// disable stack trace for the example
	check.PrintStack = true
}

// Capture the output of the go program while running function f
//
// Source https://gist.github.com/mindscratch/0faa78bd3c0005d080bf
func CaptureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()

	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}
