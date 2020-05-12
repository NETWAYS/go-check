package status

import (
	"fmt"
	"os"
)

const (
	// OK means everything is fine
	OK       = 0
	// Warning means there is a problem the admin should review
	Warning  = 1
	// Critical means there is a problem that requires immediate action
	Critical = 2
	// Unknown means the status can not be determined, probably due to an error or something missing
	Unknown  = 3
)

// String returns the string corresponding to a state
func String(status int) string {
	switch status {
	case OK:
		return "OK"
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	case Unknown:
	}
	return "UNKNOWN"
}

// Exit does plugin exit with specified rc and output
func Exit(rc int, output string, args ...interface{}) {
	fmt.Println(String(rc), "-", fmt.Sprintf(output, args...))
	os.Exit(rc)
}

// ExitError exists with an Unknown state while reporting the error
//
// TODO: more information about the error
func ExitError(err error) {
	Exit(Unknown, err.Error())
}
