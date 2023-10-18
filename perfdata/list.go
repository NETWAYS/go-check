package perfdata

import (
	"strings"
)

// PerfdataList can store multiple perfdata and brings a simple fmt.Stringer interface
// nolint: golint, revive
type PerfdataList []*Perfdata

// String returns string representations of all Perfdata
func (l PerfdataList) String() string {
	var out strings.Builder

	for _, p := range l {
		out.WriteString(" ")
		out.WriteString(p.String())
	}

	return strings.Trim(out.String(), " ")
}

// Add adds a Perfdata pointer to the list
func (l *PerfdataList) Add(p *Perfdata) {
	*l = append(*l, p)
}
