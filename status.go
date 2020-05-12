package check

const (
	// OK means everything is fine
	OK = 0
	// Warning means there is a problem the admin should review
	Warning = 1
	// Critical means there is a problem that requires immediate action
	Critical = 2
	// Unknown means the status can not be determined, probably due to an error or something missing
	Unknown = 3
)

// StatusText returns the string corresponding to a state
func StatusText(status int) string {
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
