package check

const (
	// OK means everything is fine
	OK       = 0
	OKString = "OK"
	// Warning means there is a problem the admin should review
	Warning       = 1
	WarningString = "WARNING"
	// Critical means there is a problem that requires immediate action
	Critical       = 2
	CriticalString = "CRITICAL"
	// Unknown means the status can not be determined, probably due to an error or something missing
	Unknown       = 3
	UnknownString = "UNKNOWN"
)

// StatusText returns the string corresponding to a state
func StatusText(status int) string {
	switch status {
	case OK:
		return OKString
	case Warning:
		return WarningString
	case Critical:
		return CriticalString
	case Unknown:
	}

	return UnknownString
}
