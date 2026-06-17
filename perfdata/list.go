package perfdata

import (
	"strings"
)

// PerfdataList can store multiple perfdata and implements the fmt.Stringer interface
// to provide formated output for the performance data
type PerfdataList []*Perfdata //nolint: revive

// String returns string representations of all Perfdata added to the list
func (l *PerfdataList) String() string {
	var out strings.Builder

	for _, p := range *l {
		pfDataString, err := p.ValidatedString()

		// Ignore perfdata points which fail to format
		if err == nil {
			out.WriteString(" ")
			out.WriteString(pfDataString)
		}
	}

	return strings.Trim(out.String(), " ")
}

// Add adds a Perfdata pointer to the list. Note that, it's not concurrency safe.
func (l *PerfdataList) Add(p *Perfdata) {
	*l = append(*l, p)
}
